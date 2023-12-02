// Package http_template middleware/http_template/http_template.go
package http_template

import (
	"github.com/ApolloMedTech/Middleware/alertManager"
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Render is a helper function to render a http_template with Pongo2
func Render(c *gin.Context, templateFile string, data pongo2.Context, alerts ...*Alert) {
	// Check if an alert is provided and add it to the context
	logrus.Debug("Rendering template: ", templateFile)
	if len(alerts) > 0 && alerts[0] != nil {
		logrus.Debug("Alert: ", alerts[0])
		data["alert"] = alerts[0]
	}

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

func NewAlert(alertType AlertType, message string) *Alert {
	logrus.Debug("Creating new alert: ", message)
	return &Alert{
		Type: alertType,
		Text: message,
	}
}

// NewRender is a helper function to render a http_template with Pongo2
func NewRender(c *gin.Context, templateFile string, data pongo2.Context, localizationStrings map[string]string) {

	alerts := alertManager.GetAlerts(c)
	if len(alerts) > 0 {
		logrus.Debug("Alert: ", alerts[0])
		data["alert"] = alerts[0]
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
