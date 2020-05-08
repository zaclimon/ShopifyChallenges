package server

import (
	"UtsuruConcept/controllers"
	"github.com/gin-gonic/gin"
)

func UtsuruRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/upload", controllers.Upload)
	router.POST("/search", controllers.Search)
	return router
}
