package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Type responsible for handling login information such as the user's email and password
type RegisterRequest struct {
	Email    string `form:"email" json:"email" xml:"email" binding:"required,email"`
	Password string `form:"password" json:"password" xml:"password" binding:"required,min=6"`
}

// Registers a new user through the image repository using the "/register" endpoint.
func Register(c *gin.Context) {
	var jsonRequest RegisterRequest
	err := c.ShouldBindJSON(&jsonRequest)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	_, err = models.CreateAndInsertNewUser(jsonRequest.Email, jsonRequest.Password, db.GetDb())

	if err != nil {
		errorCode := http.StatusInternalServerError
		if errors.Is(err, models.DbDuplicatedEmailError) {
			errorCode = http.StatusOK
		}
		showResponseError(c, errorCode, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "The user has been created successfully",
	})
}
