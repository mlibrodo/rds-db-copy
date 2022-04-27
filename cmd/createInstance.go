package cmd

import (
	"github.com/mikel-at-tatari/tatari-dev-db/aws/rds"
	"github.com/spf13/cobra"
)

// createInstance represents the copy command
var createInstance = &cobra.Command{
	Use:   "createInstance",
	Short: "Create a blank RDS instance",
	Long:  "Create an RDS instance based",
	Run: func(cmd *cobra.Command, args []string) {
		subnetGroupName := "project-rds-copy-vpc-rds-subnet-group"
		masterUser := "BtsY8qNqHx4xBb4r"
		masterPassword := "q2vSeHgY9fg4G88X"
		createInstance := rds.CreateInstance{
			InstanceName:    "UNCLAIMED",
			InstanceClass:   "db.t2.micro",
			SubnetGroupName: subnetGroupName,
			EngineVersion:   "12.10",
			MasterUser:      masterUser,
			MasterPassword:  masterPassword,
		}

		err, _ := createInstance.Exec()

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createInstance)
}
