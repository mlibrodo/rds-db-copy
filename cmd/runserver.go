package cmd

import (
	"github.com/mlibrodo/rds-db-copy/server"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var runServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the server",
	Long:  `Run the Server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(runServerCmd)
}
