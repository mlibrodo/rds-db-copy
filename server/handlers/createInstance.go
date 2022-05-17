package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/mlibrodo/rds-db-copy/server/middlewares"
	"github.com/mlibrodo/rds-db-copy/server/templateHelpers"
	"net/http"
	"strconv"
	"strings"
)

type CreateInstanceForm struct {
	InstanceClass string `form:"instanceClass"`
	SubnetGroup   string `form:"subnetGroup"`
	EngineVersion string `form:"engineVersion"`
	StorageSize   int32  `form:"storageSize"`
}

func POSTCreateInstance(c *gin.Context) {

	var req CreateInstanceForm
	var err error
	var inst *provisioner.RDSInstance
	var tx pgx.Tx

	if c.ShouldBind(&req) != nil {
		log.Error("Invalid input")
		c.Error(errors.New("Invalid input"))
		c.Status(400)
		return
	}
	tx = c.MustGet(middlewares.DBRequestConn).(pgx.Tx)

	engine, version := getEngineVersion(req.EngineVersion)

	inst, err = provisioner.NewRDSInstance(
		context.TODO(),
		tx,
		req.InstanceClass,
		req.SubnetGroup,
		engine,
		version,
		config.GetConfig().AWS.RDS.MasterUsername,
		config.GetConfig().AWS.RDS.MasterPassword,
		req.StorageSize)

	if err != nil {
		log.Error("Unable to create instance", err.Error())
		c.Error(errors.New("Unable to create instance"))
		c.Status(400)
		return
	}

	log.Info("Created instance", inst)
	GETCreateInstanceForm(c)
}

func GETCreateInstanceForm(c *gin.Context) {

	rdsOptions := config.GetConfig().AWS.RDS

	c.HTML(http.StatusOK, "createInstance.tmpl", gin.H{
		"InstanceClassOptions": getInstanceClassOpts(rdsOptions.SupportedInstanceClasses),
		"SubnetGroupOptions":   getSubnetOpts(rdsOptions.SubnetGroupNames),
		"EngineVersionOptions": getSupportedEngineOptions(rdsOptions.SupportedEngines),
		"StorageSizeOptions":   getStorageOptions(rdsOptions.AllowedStorageSizeGBRange),
	})
}

func getInstanceClassOpts(instanceClasses []string) []templateHelpers.SelectOption {
	var opts []templateHelpers.SelectOption
	for _, instanceClass := range instanceClasses {
		opts = append(opts, templateHelpers.SelectOption{
			Name:  instanceClass,
			Value: instanceClass,
		})
	}

	return opts
}

func getSubnetOpts(subnets []string) []templateHelpers.SelectOption {
	var opts []templateHelpers.SelectOption
	for _, subnet := range subnets {
		opts = append(opts, templateHelpers.SelectOption{
			Name:  subnet,
			Value: subnet,
		})
	}

	return opts
}

func getSupportedEngineOptions(engines []string) []templateHelpers.SelectOption {
	var opts []templateHelpers.SelectOption
	for _, engine := range engines {
		opts = append(opts, templateHelpers.SelectOption{
			Name:  engine,
			Value: engine,
		})
	}

	return opts
}

func getEngineVersion(engineVersion string) (engine string, version string) {
	splits := strings.Split(engineVersion, ":")
	engine = splits[0]
	version = splits[1]

	return engine, version
}

func getStorageOptions(gbRange string) []templateHelpers.SelectOption {
	minMax := strings.Split(gbRange, ":")
	min, _ := strconv.Atoi(minMax[0])
	max, _ := strconv.Atoi(minMax[1])

	var storageOpts []templateHelpers.SelectOption

	for i := min; i <= max; i++ {
		storageOpts = append(storageOpts,
			templateHelpers.SelectOption{Name: strconv.Itoa(i), Value: strconv.Itoa(i)})
	}

	return storageOpts
}
