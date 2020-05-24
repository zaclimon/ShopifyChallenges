package controllers

import (
	"UtsuruConcept/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Env struct {
	Db   models.UtsuruDataStore
	Mode string
}

const (
	// Test mode defines the router while running tests
	TestMode = "Test"
	// Server Mode defines the router when running Gin as a server
	ServerMode = "Server"
)

// UtsuruRouter registers the HTTP routes (endpoints) to a given function.
// It returns the router for the application.
func UtsuruRouter(env *Env) *gin.Engine {

	if env.Mode == TestMode {
		gin.DefaultWriter = ioutil.Discard
	}

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
