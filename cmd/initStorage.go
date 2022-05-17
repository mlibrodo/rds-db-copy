package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/spf13/cobra"
	"os"
)

var initStorageCmd = &cobra.Command{
	Use: "initstorage",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		var conn *pgx.Conn

		ctx := context.TODO()
		conn, err = pgx.Connect(ctx, config.GetConfig().DatabaseURL)
		defer conn.Close(ctx)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)

		tx, err := conn.Begin(ctx)
		provisioner.InitStorage(ctx, tx)
		tx.Commit(ctx)
	},
}

func init() {
	rootCmd.AddCommand(initStorageCmd)
}
