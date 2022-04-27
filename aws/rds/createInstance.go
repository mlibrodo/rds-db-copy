package rds

import (
	"context"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/config"
	"github.com/mikel-at-tatari/tatari-dev-db/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type CreateInstanceInput struct {
	InstanceName       string
	InstanceClass      string
	SubnetGroupName    string
	PubliclyAccessible bool
	EngineVersion      string
	MasterUser         string
	MasterPassword     string
}

func (c CreateInstanceInput) makeAWSCreateDBInstanceInput() *rds.CreateDBInstanceInput {

	engine := "postgres"
	storageSize := int32(5)

	return &rds.CreateDBInstanceInput{
		AllocatedStorage:     aws.Int32(storageSize),
		DBInstanceClass:      aws.String(c.InstanceClass),
		DBInstanceIdentifier: aws.String(c.InstanceName),
		Engine:               aws.String(engine),
		EngineVersion:        aws.String(c.EngineVersion),
		DBSubnetGroupName:    aws.String(c.SubnetGroupName),
		DeletionProtection:   aws.Bool(false),
		PubliclyAccessible:   aws.Bool(c.PubliclyAccessible),
		MasterUsername:       aws.String(c.MasterUser),
		MasterUserPassword:   aws.String(c.MasterPassword),
	}
}

type CreateInstanceOutput struct {
	DBHost string
	DBPort int32
}

func CreateInstance(in CreateInstanceInput) (*CreateInstanceOutput, error) {

	svc := rds.NewFromConfig(*config.AWSConfig)

	awsInput := in.makeAWSCreateDBInstanceInput()
	dbInstance, err := svc.CreateDBInstance(context.TODO(), awsInput)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	db := dbInstance.DBInstance
	log.WithFields(log.Fields{
		"DBHost": db.Endpoint.Address,
		"DBPort": db.Endpoint.Port,
	}).Info("Instance Created")

	out := CreateInstanceOutput{DBHost: *db.Endpoint.Address, DBPort: db.Endpoint.Port}
	return &out, nil
}
