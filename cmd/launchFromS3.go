package cmd

import (
	"encoding/json"
	"github.com/mlibrodo/rds-db-copy/aws/rds"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a backup in S3 to RDS",
	Long:  `Launch an S3 backup to a newly created RDS Postgres instance`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO figure out how to use VIPER with cobra to set these either as func param or ENV var
		subnetGroupName := "project-rds-copy-vpc-rds-subnet-group"
		masterUser := "BtsY8qNqHx4xBb4r" // TODO hopefully these username/password are injected via an ENV var or from some key/value store ie. VAULT
		masterPassword := "q2vSeHgY9fg4G88X"
		s3Bucket := "librodo-db-backups"
		s3Key := "source_db_1651260021.sql.tar.gz"
		instanceClass := provisioner.DefaultInstanceClass
		engineVersion := provisioner.DefaultEngineVersion
		dbName := "sinkDB"

		createInstance := &rds.CreateInstance{
			DBInstanceID:       provisioner.GenerateUnclaimedInstanceId(),
			InstanceClass:      instanceClass,
			SubnetGroupName:    subnetGroupName,
			PubliclyAccessible: true,
			EngineVersion:      engineVersion,
			MasterUser:         masterUser,
			MasterPassword:     masterPassword,
		}

		s3 := s3.S3Object{
			Bucket: s3Bucket,
			Key:    s3Key,
		}

		rdsDb := provisioner.RDSDB{createInstance}
		db, _ := rdsDb.Launch(s3, dbName)

		ret, _ := json.Marshal(db)
		print(string(ret))
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
