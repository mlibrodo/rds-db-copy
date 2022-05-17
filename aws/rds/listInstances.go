package rds

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/mlibrodo/rds-db-copy/aws/config"
)

const availStatus = "available"

func ListInstances() (*[]RDSInstanceDescriptor, error) {
	svc := rds.NewFromConfig(*config.AWSConfig)

	var descriptors []RDSInstanceDescriptor

	if out, err := svc.DescribeDBInstances(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		for _, inst := range out.DBInstances {

			if inst.Endpoint == nil || notAvailable(inst) {
				continue
			}
			descriptor := RDSInstanceDescriptor{
				DBHost:        *inst.Endpoint.Address,
				DBPort:        inst.Endpoint.Port,
				DBInstanceId:  *inst.DBInstanceIdentifier,
				DBInstanceARN: *inst.DBInstanceArn,
				DBIResourceId: *inst.DbiResourceId,
				Engine:        *inst.Engine,
				EngineVersion: *inst.EngineVersion,
				InstanceClass: *inst.DBInstanceClass,
				Status:        *inst.DBInstanceStatus,
			}
			descriptors = append(descriptors, descriptor)
		}

		return &descriptors, nil
	}
}

func notAvailable(instance types.DBInstance) bool {

	return *instance.DBInstanceStatus != availStatus
}
