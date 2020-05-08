package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type SearchRequest struct {
	Token string                `form:"token" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

func Search(c *gin.Context) {
	var requestBody SearchRequest
	err := c.ShouldBindWith(&requestBody, binding.FormMultipart)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	// Validate that the token is valid before continuing
	token, err := jwt.Parse(requestBody.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if models.IsValidImageExtension(requestBody.Image.Filename) {
			// Do the query to verify the images
			uploadFolder := os.Getenv("UPLOAD_FOLDER")
			imagePath := fmt.Sprintf("%s/%s", uploadFolder, requestBody.Image.Filename)

			if err = c.SaveUploadedFile(requestBody.Image, imagePath); err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}

			imageData, err := generateImageData(requestBody.Image.Filename)

			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}

			similarImages, err := getSimilarImages(imageData.ImageHash)

			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"images": similarImages,
			})
		} else {
			// Return error
		}
	}
}

func getSimilarImages(imageHash uint64) (*[]models.Image, error) {
	dbObj := db.GetDb()
	var images []models.Image
	hashThreshold, err := strconv.Atoi(os.Getenv("PHASH_THRESHOLD"))

	if err != nil {
		return nil, err
	}

	// Not sure if a raw SQL query is the right way to go but it works...
	sqlQuery := fmt.Sprintf("SELECT images.* from images, image_data WHERE images.id = image_data.image_id AND bit_count(%d ^ image_data.image_hash) <= %d", imageHash, hashThreshold)
	dbObj.Raw(sqlQuery).Find(&images)
	return &images, nil
}
