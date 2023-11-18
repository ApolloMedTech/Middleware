package error

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterErrorRoutes(router *gin.Engine) {
	// Handle 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		errorHandler(c, http.StatusNotFound, "Page not found")
	})

	// Handle 405 Method Not Allowed
	router.NoMethod(func(c *gin.Context) {
		errorHandler(c, http.StatusMethodNotAllowed, "Method not allowed")
	})

	// You can add more error handlers here
}
