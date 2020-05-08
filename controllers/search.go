package controllers

import (
	"UtsuruConcept/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"os"
)

type SearchRequest struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Search(c *gin.Context) {
	var multipartFormRequest SearchRequest
	err := c.ShouldBindWith(&multipartFormRequest, binding.FormMultipart)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	// Validate that the token is valid before continuing
	token, err := jwt.Parse(multipartFormRequest.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//userID := claims["id"].(string)
		//dbObj := db.GetDb()
		form, err := c.MultipartForm()

		if err != nil {
			showResponseError(c, http.StatusBadRequest, err)
			return
		}
		files := form.File["image"]
		if len(files) == 0 || len(files) >= 2 {
			// Return an error saying only one file can be looked at a time
			return
		}
		fileInfo := files[0]
		if models.IsValidImageExtension(fileInfo.Filename) {
			// Do the query to verify the images
		} else {
			// Return error
		}
	}
}
