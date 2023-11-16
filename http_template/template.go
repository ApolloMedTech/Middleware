// Package http_template middleware/http_template/http_template.go
package http_template

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Render is a helper function to render a http_template with Pongo2
func Render(c *gin.Context, templateFile string, data pongo2.Context) {

	template, err := pongo2.FromFile(templateFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template Error: "+err.Error())
		return
	}

	html, err := template.Execute(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template Execution Error: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
