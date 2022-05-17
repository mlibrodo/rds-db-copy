package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/server/templateHelpers"

	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/mlibrodo/rds-db-copy/server/middlewares"
	"net/http"
)

type RestoreDBForm struct {
	RDSInstance string `form:"rdsInstance"`
	S3Backup    string `form:"s3Backup"`
}

func POSTRestoreDB(c *gin.Context) {
	var req RestoreDBForm
	var dbCopy *provisioner.DBCopy
	var err error

	if c.ShouldBind(&req) != nil {
		log.Error("Invalid input")
		c.Error(errors.New("Invalid input"))
		c.Status(400)
		return
	}
	dbCopy, err = provisioner.NewDBCopy(
		c,
		c.MustGet(middlewares.DBRequestConn).(pgx.Tx),
		req.S3Backup,
		req.RDSInstance,
	)
	if err != nil {
		log.Error("Unable to create a copy of DB", err.Error())
		c.Error(errors.New("Unable to create a copy of DB"))
		c.Status(400)
		return
	}
	log.Info("DB Copied", dbCopy)
	GETRestoreDBForm(c)
}

func GETRestoreDBForm(c *gin.Context) {

	var err error
	instances, err := provisioner.GetAllInstances(c, c.MustGet(middlewares.DBRequestConn).(pgx.Tx))
	if err != nil {
		log.Error("Unable to get instances", err.Error())
		c.Error(errors.New("Unable to get RDS instances"))
		c.Status(400)
		return
	}
	backups, err := s3.List(&config.GetConfig().Backup.S3Bucket, nil)

	if err != nil {
		log.Error("Unable to get backups", err.Error())
		c.Error(errors.New("Unable to get backups"))
		c.Status(400)
		return
	}

	instanceOptions := make([]templateHelpers.SelectOption, len(instances))

	for i, inst := range instances {
		instanceOptions[i] = toInstanceOption(*inst)
	}

	c.HTML(http.StatusOK, "restore.tmpl",
		gin.H{"InstanceOptions": instanceOptions, "BackupOptions": getBackupOptions(backups)})
}

func toInstanceOption(i provisioner.RDSInstance) templateHelpers.SelectOption {
	return templateHelpers.SelectOption{
		Name:  fmt.Sprintf("%v", i.RDSID),
		Value: fmt.Sprintf("%v", i.RDSID),
	}
}

func getBackupOptions(backups []string) []templateHelpers.SelectOption {
	var opts []templateHelpers.SelectOption
	for _, backup := range backups {
		opts = append(opts, templateHelpers.SelectOption{
			Name:  backup,
			Value: backup,
		})
	}

	return opts
}
