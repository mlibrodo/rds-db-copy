package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/version"
	"net/http"
)

func GETVersion(c *gin.Context) {

	c.JSON(
		http.StatusOK,
		gin.H{
			"Build Date:": version.BuildDate,
			"Git Commit:": version.GitCommit,
			"Version:":    version.Version,
			"Go Version:": version.GoVersion,
			"OS / Arch:":  version.OsArch,
		})

}
