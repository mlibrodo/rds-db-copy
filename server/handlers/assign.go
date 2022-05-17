package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/provisioner"
	"github.com/mlibrodo/rds-db-copy/server/middlewares"
	"github.com/mlibrodo/rds-db-copy/server/templateHelpers"
	"net/http"
	"strconv"
)

type AssignDBForm struct {
	IAMUser  string `form:"iamUser"`
	DBCopyId int32  `form:"dbCopyId"`
}

func POSTAssign(c *gin.Context) {
	var req AssignDBForm
	var err error
	var dbCopy *provisioner.DBCopy

	if c.ShouldBind(&req) != nil {
		log.Error("Invalid input")
		c.Error(errors.New("Invalid input"))
		c.Status(400)
		return
	}

	dbCopy, err = provisioner.GetDBCopy(c, c.MustGet(middlewares.DBRequestConn).(pgx.Tx), req.DBCopyId)

	if err != nil {
		log.Error("Error GetDBCopy", err.Error())
		c.Error(errors.New("Error GetDBCopy"))
		c.Status(400)
		return
	}

	dbCopy, err = provisioner.AssignCopyToUser(c, c.MustGet(middlewares.DBRequestConn).(pgx.Tx), *dbCopy, req.IAMUser)
	if err != nil {
		log.Error("Error Assigning", err.Error())
		c.Error(errors.New("Error Assigning"))
		c.Status(400)
		return
	}

	GETAssignDBForm(c)
}

func GETAssignDBForm(c *gin.Context) {

	var err error
	var dbCopies []*provisioner.DBCopy
	dbCopies, err = provisioner.GetAllUnassignedCopies(c, c.MustGet(middlewares.DBRequestConn).(pgx.Tx))

	if err != nil {
		log.Error("Unable to get unassigned copies", err.Error())
		c.Error(errors.New("Unable to get unassigned copies"))
		c.Status(400)
		return
	}

	c.HTML(http.StatusOK, "assign.tmpl", gin.H{"DBCopies": getOptions(dbCopies)})
}

func getOptions(copies []*provisioner.DBCopy) []templateHelpers.SelectOption {
	var opts []templateHelpers.SelectOption
	for _, ptr := range copies {
		dbCopy := *ptr
		opts = append(opts, templateHelpers.SelectOption{
			Name:  dbCopy.DbName,
			Value: strconv.Itoa(int(dbCopy.ID)),
		})
	}
	return opts
}
