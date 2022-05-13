package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/postgres"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"net/http"
	"strconv"
	"strings"
)

type BackupDBForm struct {
	DBOption string `form:"dbOption"`
}

func PostBackupDB(c *gin.Context) {
	var req BackupDBForm

	if c.ShouldBind(&req) == nil {
		host, port, dbName := fromDBOption(req.DBOption)

		var dbUser string
		var dbPass string

		dbs := config.GetConfig().DBRegistry.DBs

		for _, db := range dbs {
			if host == db.Host && port == db.Port && dbName == db.DBName {
				dbUser = db.DBUser
				dbPass = db.DBPassword
				break
			}
		}

		db := &conn.PGConnInfo{
			DBHost:   host,
			DBPort:   port,
			DBName:   dbName,
			Username: dbUser,
			Password: dbPass,
		}
		s3Obj := s3.S3Object{
			Bucket: config.GetConfig().Backup.S3Bucket,
			Key:    postgres.MakeS3Key(db.DBName),
		}

		backup := postgres.BackupToS3{DB: db}
		if err := backup.Exec(s3Obj); err == nil {
			GetBackupDBForm(c)
		} else {
			log.Error("Unable to backup to s3")
			c.Status(400)
		}
	} else {
		log.Error("Invalid input")
		c.Status(400)
	}

}

func GetBackupDBForm(c *gin.Context) {

	dbs := config.GetConfig().DBRegistry.DBs

	dbOptions := make([]string, len(dbs))

	for i, db := range dbs {
		dbOption := toDBOption(db)
		dbOptions[i] = dbOption
	}

	c.HTML(http.StatusOK, "backup.tmpl", gin.H{"DBOptions": dbOptions})

}

func toDBOption(db config.DB) string {
	return fmt.Sprintf("%v:%v:%v", db.Host, db.Port, db.DBName)
}

func fromDBOption(dbId string) (dbHost string, port int32, dbName string) {
	splits := strings.Split(dbId, ":")
	dbHost = splits[0]
	p, _ := strconv.Atoi(splits[1])
	port = int32(p)
	dbName = splits[2]

	return
}
