package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/mlibrodo/rds-db-copy/config"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

var templates *template.Template

func Start() {

	engine := gin.Default()

	addRoutes(engine)
	templates = template.Must(template.ParseFS(templateFS, "templates/*.tmpl"))
	engine.SetHTMLTemplate(templates)

	engine.Run(config.GetConfig().Server.BindAddress)
}
