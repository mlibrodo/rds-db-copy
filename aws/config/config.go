package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var AWSConfig *aws.Config

func init() {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("personal"))

	if err != nil {
		panic(err)
	}

	AWSConfig = &cfg
}
