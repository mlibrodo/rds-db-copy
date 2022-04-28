package cmd

import (
	"github.com/mikel-at-tatari/tatari-dev-db/aws/iam"
	"github.com/spf13/cobra"
)

// doStuff represents the version command
var doStuff = &cobra.Command{
	Use:   "dostuff",
	Short: "dostuff",
	Long:  `All software has versions. This is generated code example`,
	Run: func(cmd *cobra.Command, args []string) {
		tVars := iam.TemplateVars{
			Region:        "us-east-1",
			AccountID:     "283492314377",
			DbiResourceId: "db-SLPFMIWCNA537LFX2PWLI73GSU",
			DBUserName:    "mike-iam-admin",

			AWSUser:        "mike-iam-admin",
			DBInstanceName: "test-mike9",
		}

		if err := tVars.AttachPolicyToUser(); err != nil {
			print(err.Error())
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(doStuff)
}
