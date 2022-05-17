package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/log"
	"os"
)

func DBSessionMiddleware(c *gin.Context) {
	var err error
	var conn *pgx.Conn
	var tx pgx.Tx

	// TODO use a connection pool instead of new connections to DB
	conn, err = pgx.Connect(c, config.GetConfig().DatabaseURL)

	defer conn.Close(c)

	if err != nil {
		log.Error(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	tx, err = conn.Begin(c)
	c.Set(DBRequestConn, tx)

	// Pass on to the next-in-chain
	c.Next()

	// rollback if there are any errors
	for _, err = range c.Errors {
		tx.Rollback(c)
		return
	}

	tx.Commit(c)
}
