package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"os"
	"strconv"
)

type SearchRequest struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Search(c *gin.Context) {
	var multipartFormRequest SearchRequest
	err := c.ShouldBindWith(&multipartFormRequest, binding.FormMultipart)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	// Validate that the token is valid before continuing
	token, err := jwt.Parse(multipartFormRequest.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//userID := claims["id"].(string)
		dbObj := db.GetDb()
		form, err := c.MultipartForm()

		if err != nil {
			showResponseError(c, http.StatusBadRequest, err)
			return
		}
		files := form.File["image"]
		if len(files) == 0 || len(files) >= 2 {
			// Return an error saying only one file can be looked at a time
			return
		}
		fileInfo := files[0]
		if models.IsValidImageExtension(fileInfo.Filename) {
			// Do the query to verify the images
			uploadFolder := os.Getenv("UPLOAD_FOLDER")
			imagePath := fmt.Sprintf("%s/%s", uploadFolder, fileInfo.Filename)
			err = c.SaveUploadedFile(fileInfo, imagePath)
			defer os.Remove(imagePath)
			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}
			imageObj, err := models.DecodeImage(imagePath)
			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}
			imageData, err := models.CreateImageData(imageObj)
			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}
			// Do an SQL request to verify the image based on its ImageHash
			fmt.Printf("PHash of image: %d\n", imageData.ImageHash)
			var images []models.Image
			hashThreshold, err := strconv.Atoi(os.Getenv("PHASH_THRESHOLD"))

			if err != nil {
				showResponseError(c, http.StatusInternalServerError, err)
				return
			}

			sqlQuery := fmt.Sprintf("SELECT images.* from images, image_data WHERE images.id = image_data.image_id AND bit_count(%d ^ image_data.image_hash) <= %d", imageData.ImageHash, hashThreshold)
			dbObj.Raw(sqlQuery).Find(&images)
			c.JSON(http.StatusOK, gin.H{
				"images": images,
			})
		} else {
			// Return error
		}
	}
}
