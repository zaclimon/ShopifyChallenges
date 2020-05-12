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
	router.POST("/register", env.Register)
	router.POST("/login", env.Login)
	router.POST("/upload", env.Upload)
	router.POST("/search", env.Search)
	return router
}
