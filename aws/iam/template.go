package iam

import (
	"embed"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "templates/*.tmpl")

	if err != nil {
		panic(err)
	}
}
