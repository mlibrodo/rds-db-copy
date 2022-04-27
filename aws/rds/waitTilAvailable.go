package rds

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/config"
	"time"
)

func WaitTilAvailable(instanceId *string) error {
	svc := rds.NewFromConfig(*config.AWSConfig)

	waiter := rds.NewDBInstanceAvailableWaiter(svc)
	instance := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: instanceId,
	}
	if err := waiter.Wait(context.TODO(), instance, 10*time.Minute); err != nil {
		return err
	}
	return nil
}
