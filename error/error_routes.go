package error

import (
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterErrorRoutes(router *gin.Engine, cfg config.Config) {
	// Handle 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		errorHandler(c, http.StatusNotFound, "Page not found", cfg)
	})

	// Handle 405 Method Not Allowed
	router.NoMethod(func(c *gin.Context) {
		errorHandler(c, http.StatusMethodNotAllowed, "Method not allowed", cfg)
	})

	// You can add more error handlers here
}
