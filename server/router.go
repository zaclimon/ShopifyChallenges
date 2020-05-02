package server

import (
	"UtsuruConcept/controllers"
	"github.com/gin-gonic/gin"
)

func UtsuruRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	return router
}
