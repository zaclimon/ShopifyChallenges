package controllers

import (
	"UtsuruConcept/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Type responsible for handling login information such as the user's email and password
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login a user when he/she goes through the "/login" endpoint
func (env *Env) Login(c *gin.Context) {
	var jsonRequest RegisterRequest
	err := c.ShouldBindJSON(&jsonRequest)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	user, err := env.db.GetUserByEmail(jsonRequest.Email)

	if err != nil {
		showResponseError(c, http.StatusNotFound, err)
		return
	}

	if err = models.ValidatePassword(user.Password, jsonRequest.Password); err != nil {
		showResponseError(c, http.StatusForbidden, models.InvalidCredentialsError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": user.ID,
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		showResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
	})
}
