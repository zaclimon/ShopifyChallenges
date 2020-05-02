package server

import (
	"UtsuruConcept/controllers"
	"github.com/gin-gonic/gin"
)

func UtsuruRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/register", controllers.Register)
	return router
}
