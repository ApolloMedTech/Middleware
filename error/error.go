package error

import (
	"fmt"
	"github.com/ApolloMedTech/Middleware/http_template"
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func CustomErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic
				logrus.Panic("Panic recovered: ", r)
				errorHandler(c, http.StatusInternalServerError, fmt.Sprintf("%s", r))
			}
		}()
		c.Next() // Process request
		// Check if there is an error after processing the request
		if len(c.Errors) > 0 {
			// Handle the first error
			err := c.Errors[0].Err
			logrus.Error("Error: ", err.Error())
			errorHandler(c, http.StatusInternalServerError, err.Error())
		}
	}
}

func errorHandler(c *gin.Context, statusCode int, message string) {
	// Set the status code and render the template
	http_template.Render(c, "templates/error.html", pongo2.Context{"status": statusCode, "message": message})
}
