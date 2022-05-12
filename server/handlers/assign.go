package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/log"
	"net/http"
)

type AssignDBForm struct {
	IAMUser       string `form:"iamUser"`
	RDSInstanceId string `form:"rdsInstanceId"`
	DBName        string `form:"dbName"`
	RDSRegion     string `form:"rdsRegion"`
	AWSAccountID  string `form:"awsAccountId"`
}

func Assign(c *gin.Context) {
	var req AssignDBForm
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&req) == nil {
		log.Print(req)
	}

	c.HTML(http.StatusOK, "index.tmpl", nil)
}
