package rds

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/mlibrodo/rds-db-copy/aws/config"
)

type DescribeInstanceOutput struct {
	DBHost *string
	DBPort int32
}

func DescribeInstance(instanceId *string) (*RDSInstanceDescriptor, error) {
	svc := rds.NewFromConfig(*config.AWSConfig)

	instanceFilter := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: instanceId,
	}

	if out, err := svc.DescribeDBInstances(context.TODO(), instanceFilter); err != nil {
		return nil, err
	} else {
		descriptor := RDSInstanceDescriptor{
			DBHost:        *out.DBInstances[0].Endpoint.Address,
			DBPort:        out.DBInstances[0].Endpoint.Port,
			DBInstanceId:  *out.DBInstances[0].DBInstanceIdentifier,
			DBInstanceARN: *out.DBInstances[0].DBInstanceArn,
			DBIResourceId: *out.DBInstances[0].DbiResourceId,
			Engine:        *out.DBInstances[0].Engine,
			EngineVersion: *out.DBInstances[0].EngineVersion,
			InstanceClass: *out.DBInstances[0].DBInstanceClass,
			Status:        *out.DBInstances[0].DBInstanceStatus,
		}

		return &descriptor, nil
	}

}
