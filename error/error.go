package error

import (
	"fmt"
	"net/http"

	"github.com/ApolloMedTech/Middleware/config"
	"github.com/ApolloMedTech/Middleware/http_template"
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CustomErrorHandling(cfg config.TemplatesConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic
				logrus.Panic("Panic recovered: ", r)
				errorHandler(c, http.StatusInternalServerError, fmt.Sprintf("%s", r), cfg)
			}
		}()
		c.Next() // Process request
		// Check if there is an error after processing the request
		if len(c.Errors) > 0 {
			// Handle the first error
			err := c.Errors[0].Err
			logrus.Error("Error: ", err.Error())
			errorHandler(c, http.StatusInternalServerError, err.Error(), cfg)
		}
	}
}

func errorHandler(c *gin.Context, statusCode int, message string, config config.TemplatesConfig) {
	// Set the status code and render the template
	http_template.Render(c, config.Path+"/error.html", pongo2.Context{"status": statusCode, "message": message})
}
