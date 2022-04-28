package rds

import (
	"context"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/config"
	"github.com/mikel-at-tatari/tatari-dev-db/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type CreateInstance struct {
	InstanceName       string
	InstanceClass      string
	SubnetGroupName    string
	PubliclyAccessible bool
	EngineVersion      string
	MasterUser         string
	MasterPassword     string
}

func (ins *CreateInstance) makeAWSCreateDBInstanceInput() *rds.CreateDBInstanceInput {

	engine := "postgres"
	storageSize := int32(5)
	backupRentention := int32(0) // no backups

	return &rds.CreateDBInstanceInput{
		AllocatedStorage:                aws.Int32(storageSize),
		BackupRetentionPeriod:           aws.Int32(backupRentention),
		DBInstanceClass:                 aws.String(ins.InstanceClass),
		DBInstanceIdentifier:            aws.String(ins.InstanceName),
		Engine:                          aws.String(engine),
		EngineVersion:                   aws.String(ins.EngineVersion),
		DBSubnetGroupName:               aws.String(ins.SubnetGroupName),
		DeletionProtection:              aws.Bool(false),
		PubliclyAccessible:              aws.Bool(ins.PubliclyAccessible),
		EnablePerformanceInsights:       aws.Bool(false),
		MasterUsername:                  aws.String(ins.MasterUser),
		MasterUserPassword:              aws.String(ins.MasterPassword),
		EnableIAMDatabaseAuthentication: aws.Bool(true),
	}
}

func (ins *CreateInstance) Exec() (*RDSInstanceDescriptor, error) {

	svc := rds.NewFromConfig(*config.AWSConfig)

	awsInput := ins.makeAWSCreateDBInstanceInput()
	dbInstance, err := svc.CreateDBInstance(context.TODO(), awsInput)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	db := dbInstance.DBInstance

	var host string
	var port int32
	if db.Endpoint != nil {
		host = *db.Endpoint.Address
		port = db.Endpoint.Port
	}

	log.WithFields(log.Fields{
		"DBHost":        host,
		"DBPort":        port,
		"DBInstanceId":  db.DBInstanceIdentifier,
		"DBInstanceARN": db.DBInstanceArn,
	}).Info("RDSInstanceDescriptor Created")

	out := RDSInstanceDescriptor{
		DBHost:        host,
		DBPort:        port,
		DBInstanceId:  db.DBInstanceIdentifier,
		DBInstanceARN: db.DBInstanceArn,
		DBIResourceId: db.DbiResourceId,
	}
	return &out, nil
}
