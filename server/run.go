package server

import (
	"embed"
	"github.com/gin-gonic/gin"
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

	engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
