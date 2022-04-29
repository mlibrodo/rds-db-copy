package cmd

import (
	"encoding/json"
	"github.com/mlibrodo/rds-db-copy/aws"
	"github.com/mlibrodo/rds-db-copy/aws/config"
	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/spf13/cobra"
)

var assignDBToUserCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign an RDS Postgres instance to a IAM user",
	Long: `Assign an RDS Postgres instance to an IAM user. 
This will allow the user access the instance using their IAM credentials.
For more information see: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.AWSCLI.PostgreSQL.html

To login to the assigned instance:
export PGPASSWORD="$(aws rds generate-db-auth-token --hostname <DBHOST> --port <DBPORT> --region <REGION> --username <DBUSER> )"
PGSSLMODE=verify-full PGSSLROOTCERT=<PATH_TO_AWS_RDS_CERT> psql --host=<DBHOST> --port=<DBPORT>  --dbname=<DBNAME> --username=<DBUSER>
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO figure out how to use VIPER with cobra to set these either as func param or ENV var
		region := config.AWSConfig.Region
		accountId := aws.AWSAccountId
		dbInstanceId := "unclaimed-1651260226"
		iamUser := "mike-iam-admin"
		masterUser := "BtsY8qNqHx4xBb4r" // TODO hopefully these username/password are injected via an ENV var or from some key/value store ie. VAULT
		masterPassword := "q2vSeHgY9fg4G88X"
		dbName := "sinkDB"

		assignToUser := provisioner.AssignRDSToIAMUser{
			IAMUser:          iamUser,
			DBInstanceID:     dbInstanceId,
			DBName:           dbName,
			DBMasterUser:     masterUser,
			DBMasterPassword: masterPassword,
			AWSRegionOfDB:    region,
			AWSAccountIdOfDB: accountId,
		}

		db, dbUser, _ := assignToUser.AssignToIAMUser()

		retVal := map[string]any{
			"DBInfo":       *db,
			"DBUser":       dbUser,
			"AWSAccountId": accountId,
		}

		ret, _ := json.Marshal(retVal)
		print(string(ret))
	},
}

func init() {
	rootCmd.AddCommand(assignDBToUserCmd)
}
