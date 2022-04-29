package cmd

import (
	"encoding/json"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/postgres"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a Postgres Database to S3",
	Long:  `Backup a Postgres Database to S3`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO figure out how to use VIPER with cobra to set these either as func param or ENV var
		dbHost := "localhost"
		dbPort := 5432
		dbName := "source_db"
		dbUser := "miketatari"
		dbUserPass := ""
		bucket := "librodo-db-backups"

		db := &conn.PGConnInfo{
			DBHost:   dbHost,
			DBPort:   int32(dbPort),
			DBName:   dbName,
			Username: dbUser,
			Password: dbUserPass,
		}
		s3 := s3.S3Object{
			Bucket: bucket,
			Key:    postgres.MakeS3Key(db.DBName),
		}

		backup := postgres.BackupToS3{DB: db}
		backup.Exec(s3)

		ret, _ := json.Marshal(s3)
		println(string(ret))
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
