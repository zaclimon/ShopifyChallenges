package controllers

import (
	"UtsuruConcept/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Type responsible for handling login information such as the user's email and password
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Registers a new user through the image repository using the "/register" endpoint.
func (env *Env) Register(c *gin.Context) {
	var jsonRequest RegisterRequest
	err := c.ShouldBindJSON(&jsonRequest)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	user, err := models.CreateNewUser(jsonRequest.Email, jsonRequest.Password)
	if err != nil {
		showResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if err = env.Db.InsertOrUpdateUser(user); err != nil {
		showResponseError(c, http.StatusConflict, models.DbDuplicatedEmailError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "The user has been created successfully",
	})
}
