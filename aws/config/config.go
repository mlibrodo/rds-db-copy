package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsSvcCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/mlibrodo/rds-db-copy/config"
)

var AWSConfig *aws.Config

func init() {

	var cfg aws.Config
	var err error
	if profile := config.GetConfig().AWS.ServiceConfigProfile; profile != "" {
		cfg, err = awsSvcCfg.LoadDefaultConfig(context.TODO(),
			awsSvcCfg.WithSharedConfigProfile(profile))
	} else {
		cfg, err = awsSvcCfg.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		panic(err)
	}

	AWSConfig = &cfg
}
