package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/server/middlewares"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

var templates *template.Template

func Start() {

	engine := gin.Default()
	engine.Use(middlewares.DBSessionMiddleware)
	templates = template.Must(template.ParseFS(templateFS, "templates/*.tmpl"))
	engine.SetHTMLTemplate(templates)
	addRoutes(engine)
	engine.Run(config.GetConfig().Server.BindAddress)
}
