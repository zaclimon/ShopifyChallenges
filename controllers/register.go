package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Email    string `form:"email" json:"email" xml:"email" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
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

	user, err := models.CreateAndInsertNewUser(jsonRequest.Email, jsonRequest.Password, db.GetDb())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
