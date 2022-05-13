package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/log"
	"net/http"
)

type RestoreDBForm struct {
	IAMUser       string `form:"iamUser"`
	RDSInstanceId string `form:"rdsInstanceId"`
	DBName        string `form:"dbName"`
	RDSRegion     string `form:"rdsRegion"`
	AWSAccountID  string `form:"awsAccountId"`
}

func POSTLaunch(c *gin.Context) {
	var req RestoreDBForm

	if c.ShouldBind(&req) == nil {
		log.Print(req)
	}

	c.HTML(http.StatusOK, "assign.tmpl", nil)
}
