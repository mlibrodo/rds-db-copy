package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/server/handlers"
	"net/http"
)

func addRoutes(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)

	})

	engine.POST("/assign", func(c *gin.Context) {
		handlers.Assign(c)
	})
}
