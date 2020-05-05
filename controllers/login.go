package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Type responsible for handling login information
// such as the user's email and password
type LoginRequest struct {
	Email    string `form:"user" json:"user" xml:"user" binding:"required,email"`
	Password string `form:"user" json:"user" xml:"user" binding:"required,min=6"`
}

func Login(c *gin.Context) {
	var jsonRequest RegisterRequest
	err := c.ShouldBindJSON(&jsonRequest)
	dbObj := db.GetDb()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := models.GetUserByEmail(jsonRequest.Email, jsonRequest.Password, dbObj)

	if err != nil {
		errorCode := http.StatusInternalServerError
		if errors.Is(err, models.InvalidCredentialsError) {
			errorCode = http.StatusForbidden
		} else if errors.Is(err, models.DbUserNotFoundError) {
			errorCode = http.StatusNotFound
		}
		c.JSON(errorCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": user.ID,
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}
