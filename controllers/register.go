package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Email    string `form:"email" json:"email" xml:"email" binding:"required,email"`
	Password string `form:"password" json:"password" xml:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var jsonRequest RegisterRequest
	err := c.ShouldBindJSON(&jsonRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = models.CreateAndInsertNewUser(jsonRequest.Email, jsonRequest.Password, db.GetDb())

	if err != nil {
		errorCode := http.StatusInternalServerError
		if errors.Is(err, models.DbDuplicatedEmailError) {
			errorCode = http.StatusOK
		}
		c.JSON(errorCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "The user has been created successfully",
	})
}
