// Package templateManager middleware/http_template/http_template.go
package templateManager

import (
	"github.com/ApolloMedTech/Middleware/alertManager"
	"github.com/ApolloMedTech/Middleware/localization"
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Render is a helper function to render a http_template with Pongo2
func Render(c *gin.Context, templateFile string, data pongo2.Context, localizationStrings map[string]string) {

	// Check if an alert is provided and add it to the context
	alerts := alertManager.GetAlerts(c)
	if len(alerts) > 0 {
		logrus.Debug("Alert: ", alerts[0])
		data["alert"] = alerts[0]
	}

	// Add the base strings to the localizationStrings map
	baseStrings := localization.LocalizePrefixStrings(c, "base_")

	for k, v := range baseStrings {
		localizationStrings[k] = v
	}

	data["localizedStrings"] = localizationStrings

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
	alertManager.ClearAlerts(c)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
