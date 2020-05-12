package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
)

// showResponseError writes to the response body the error retrieved when executing functions
func showResponseError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

// savedUploadedFile persists the file in the local filesystem for further processing.
// It returns an error if the file could not be saved.
func saveUploadedFile(fileInfo *multipart.FileHeader, c *gin.Context) error {
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	destinationPath := fmt.Sprintf("%s/%s", uploadFolder, fileInfo.Filename)
	if err := c.SaveUploadedFile(fileInfo, destinationPath); err != nil {
		return err
	}
	return nil
}

// validateToken verifies the validity of a token when doing an authenticated request.
// It returns an error if the token could not be parsed correctly.
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
