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
	Engine             string
	EngineVersion      string
	// MasterUser is the initial user to bootstrap the PG RDS instance.
	// Only alphanumeric usernames are allowed.
	// For more information visit: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.MasterAccounts.html
	// TODO check for alphanumeric master username
	MasterUser string
	// MasterPassword password for the MasterUser
	MasterPassword string
	StorageSize    int32
}

func (ins *CreateInstance) makeAWSCreateDBInstanceInput() *rds.CreateDBInstanceInput {

	return &rds.CreateDBInstanceInput{
		AllocatedStorage:                aws.Int32(ins.StorageSize),
		DBInstanceClass:                 aws.String(ins.InstanceClass),
		DBInstanceIdentifier:            aws.String(ins.DBInstanceID),
		Engine:                          aws.String(ins.Engine),
		EngineVersion:                   aws.String(ins.EngineVersion),
		DBSubnetGroupName:               aws.String(ins.SubnetGroupName),
		PubliclyAccessible:              aws.Bool(ins.PubliclyAccessible),
		MasterUsername:                  aws.String(ins.MasterUser),
		MasterUserPassword:              aws.String(ins.MasterPassword),
		EnablePerformanceInsights:       aws.Bool(false),
		DeletionProtection:              aws.Bool(false),
		EnableIAMDatabaseAuthentication: aws.Bool(true),
		BackupRetentionPeriod:           aws.Int32(int32(0)), // no backups
	}
}

func (ins *CreateInstance) Exec() (*RDSInstanceDescriptor, error) {

	svc := rds.NewFromConfig(*config.AWSConfig)

	awsInput := ins.makeAWSCreateDBInstanceInput()
	resp, err := svc.CreateDBInstance(context.TODO(), awsInput)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	db := resp.DBInstance

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
		DBInstanceId:  *db.DBInstanceIdentifier,
		DBInstanceARN: *db.DBInstanceArn,
		DBIResourceId: *db.DbiResourceId,
		Engine:        *db.Engine,
		EngineVersion: *db.EngineVersion,
		InstanceClass: *db.DBInstanceClass,
		Status:        *db.DBInstanceStatus,
	}
	return &out, nil
}
