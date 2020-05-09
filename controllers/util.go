package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"strconv"
)

// showResponseError writes to the response body the error retrieved when executing functions
func showResponseError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

// savedUploadedFile persists the file in the local filesystem for further processing.
func saveUploadedFile(fileInfo *multipart.FileHeader, c *gin.Context) error {
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	destinationPath := fmt.Sprintf("%s/%s", uploadFolder, fileInfo.Filename)
	if err := c.SaveUploadedFile(fileInfo, destinationPath); err != nil {
		return err
	}
	return nil
}

// generateImageData creates image metadata that can be used for further processing.
func generateImageData(filePath string) (*models.ImageData, error) {
	decodedImage, err := models.DecodeImage(filePath)

	if err != nil {
		return nil, err
	}

	imageData, err := models.CreateImageData(decodedImage)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

// validateToken verifies the validity of a token when doing an authenticated request.
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

// getSimilarImages retrieves images that are "similar" based on a given image hash.
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
