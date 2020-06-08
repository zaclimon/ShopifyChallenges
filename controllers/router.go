package controllers

import (
	"UtsuruConcept/models"
	"UtsuruConcept/storage"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Env struct {
	Db   models.UtsuruDataStore
	Gcs  storage.UtsuruFileStorage
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
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello world!")
	})

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

// savedUploadedFile persists the file in the local filesystem for further processing.
// It returns an error if the file could not be saved.
func saveUploadedFile(fileInfo *multipart.FileHeader, c *gin.Context) error {
	destinationPath := fmt.Sprintf("%s/%s", os.TempDir(), fileInfo.Filename)
	if err := c.SaveUploadedFile(fileInfo, destinationPath); err != nil {
		return err
	}
	return nil
}

// validateToken verifies the validity of a token when doing an authenticated request.
// It returns an error if the token could not be parsed correctly.
func validateToken(token string) (string, error) {
	// Validate that the token is valid before continuing
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["id"].(string), nil
	}

	return "", err
}
