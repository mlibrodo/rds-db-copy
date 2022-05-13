package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/server/handlers"
	"net/http"
)

func addRoutes(engine *gin.Engine) {
	engine.GET("/version", func(c *gin.Context) {
		handlers.Version(c)
	})

	engine.GET("/assign", func(c *gin.Context) {
		c.HTML(http.StatusOK, "assign.tmpl", nil)
	})
	engine.POST("/assign", func(c *gin.Context) {
		handlers.Assign(c)
	})

	engine.POST("/launch", func(c *gin.Context) {
		handlers.PostLaunch(c)
	})

	engine.GET("/backup", func(c *gin.Context) {
		handlers.GetBackupDBForm(c)
	})
	engine.POST("/backup", func(c *gin.Context) {
		handlers.PostBackupDB(c)
	})
}
