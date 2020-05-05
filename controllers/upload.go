package controllers

import (
	"UtsuruConcept/models"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type UploadRequest struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Upload(c *gin.Context) {
	var jsonRequest UploadRequest
	err := c.ShouldBindJSON(&jsonRequest)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	// Validate that the token is valid before continuing
	token, err := jwt.Parse(jsonRequest.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate that the file is an image.
		// Let's obtain the files first
		form, err := c.MultipartForm()
		fmt.Println(claims["id"])

		if err != nil {
			showResponseError(c, http.StatusBadRequest, err)
			return
		}

		files := form.File["images"]
		for _, file := range files {
			// Verify the extension of the file
			if models.IsValidImageExtension(file.Filename) {
				// Upload to cloud storage
				// 1. Ensure that the bucket (folder) for the user exists. Create it if not available
				// 2. Verify if file metadata doesn't exist first.
				// 3. Create the file metadata in the database.
				// 4. Upload the file to the bucket (If filename already exists, don't upload it twice --> Filename based upload instead of ID based)
				userID := claims["id"].(string)
				ctx := context.Background()
				client, err := storage.NewClient(ctx)
				if err != nil {
					// The client has not been configured correctly
					showResponseError(c, http.StatusInternalServerError, err)
					return
				}

				userBucket := client.Bucket(fmt.Sprintf("images/%s", userID))
				_, err = userBucket.Attrs(ctx)

				if err != nil {
					// Bucket does not exist yet
					projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
					if err = userBucket.Create(ctx, projectID, nil); err != nil {
						// Error happened while trying to create the user bucket.
						showResponseError(c, http.StatusInternalServerError, err)
					}
				}
			}
		}

	} else {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}
}

func showResponseError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
