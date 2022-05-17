package aws

import "github.com/mlibrodo/rds-db-copy/config"

var AWSAccountId string

func init() {
	AWSAccountId = config.GetConfig().AWS.AccountId
}
