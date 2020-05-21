package controllers

import (
	"UtsuruConcept/models"
	"github.com/gin-gonic/gin"
)

type Env struct {
	db models.UtsuruDataStore
}

// UtsuruRouter registers the HTTP routes (endpoints) to a given function.
// It returns the router for the application.
func UtsuruRouter(db *models.DB) *gin.Engine {

	env := &Env{db}

	router := gin.Default()
	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.POST("/register", env.Register)
			v1.POST("/login", env.Login)
			v1.POST("/upload", env.Upload)
			v1.POST("/search", env.Search)
		}
	}

	return router
}
