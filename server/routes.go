package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/server/handlers"
	"net/http"
)

func addRoutes(engine *gin.Engine) {
	engine.GET("/version", func(c *gin.Context) {
		handlers.GETVersion(c)
	})

	engine.GET("/assign", func(c *gin.Context) {
		c.HTML(http.StatusOK, "assign.tmpl", nil)
	})
	engine.POST("/assign", func(c *gin.Context) {
		handlers.POSTAssign(c)
	})

	engine.POST("/launch", func(c *gin.Context) {
		handlers.POSTLaunch(c)
	})

	engine.GET("/backup", func(c *gin.Context) {
		handlers.GETBackupDBForm(c)
	})
	engine.POST("/backup", func(c *gin.Context) {
		handlers.POSTBackupDB(c)
	})
}
