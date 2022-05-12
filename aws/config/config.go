package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var AWSConfig *aws.Config

const config_profile = "personal"

func init() {

	var cfg aws.Config
	var err error
	if config_profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(config_profile))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		panic(err)
	}

	AWSConfig = &cfg
}
