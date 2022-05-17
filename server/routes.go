package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/server/handlers"
)

func addRoutes(engine *gin.Engine) {
	engine.GET("/version", func(c *gin.Context) {
		handlers.GETVersion(c)
	})

	engine.GET("/backup", func(c *gin.Context) {
		handlers.GETBackupDBForm(c)
	})
	engine.POST("/backup", func(c *gin.Context) {
		handlers.POSTBackupDB(c)
	})

	engine.GET("/createInstance", func(c *gin.Context) {
		handlers.GETCreateInstanceForm(c)
	})
	engine.POST("/createInstance", func(c *gin.Context) {
		handlers.POSTCreateInstance(c)
	})

	engine.GET("/assign", func(c *gin.Context) {
		handlers.GETAssignDBForm(c)
	})
	engine.POST("/assign", func(c *gin.Context) {
		handlers.POSTAssign(c)
	})

	engine.GET("/restore", func(c *gin.Context) {
		handlers.GETRestoreDBForm(c)
	})
	engine.POST("/restore", func(c *gin.Context) {
		handlers.POSTRestoreDB(c)
	})

}
