package server

import (
	"UtsuruConcept/controllers"
	"github.com/gin-gonic/gin"
)

// UtsuruRouter registers the HTTP routes (endpoints) to a given function.
// It returns the router for the application.
func UtsuruRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/upload", controllers.Upload)
	router.POST("/search", controllers.Search)
	return router
}
