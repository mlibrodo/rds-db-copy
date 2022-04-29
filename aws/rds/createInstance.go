package rds

import (
	"context"
	"github.com/mlibrodo/rds-db-copy/aws/config"
	"github.com/mlibrodo/rds-db-copy/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type CreateInstance struct {
	DBInstanceID       string
	InstanceClass      string
	SubnetGroupName    string
	PubliclyAccessible bool
	EngineVersion      string
	// MasterUser is the initial user to bootstrap the PG RDS instance.
	// Only alphanumeric usernames are allowed.
	// For more information visit: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.MasterAccounts.html
	// TODO check for alphanumeric master username
	MasterUser string
	// MasterPassword password for the MasterUser
	MasterPassword string
}

func (ins *CreateInstance) makeAWSCreateDBInstanceInput() *rds.CreateDBInstanceInput {

	engine := "postgres"
	storageSize := int32(5)
	backupRentention := int32(0) // no backups

	return &rds.CreateDBInstanceInput{
		AllocatedStorage:                aws.Int32(storageSize),
		BackupRetentionPeriod:           aws.Int32(backupRentention),
		DBInstanceClass:                 aws.String(ins.InstanceClass),
		DBInstanceIdentifier:            aws.String(ins.DBInstanceID),
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
		"DBInstanceId":  *db.DBInstanceIdentifier,
		"DBInstanceARN": *db.DBInstanceArn,
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
