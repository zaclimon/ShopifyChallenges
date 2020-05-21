package controllers

import (
	"UtsuruConcept/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mime/multipart"
	"net/http"
	"os"
)

// SearchRequest handles information for the lookup of images based on a similar one.
type SearchRequest struct {
	AccessToken 	string                `form:"access_token" binding:"required"`
	Image       	*multipart.FileHeader `form:"image" binding:"required"`
}

// Search retrieves similar images from requests made on the "/search" endpoint.
func (env *Env) Search(c *gin.Context) {
	var requestBody SearchRequest
	err := c.ShouldBindWith(&requestBody, binding.FormMultipart)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	_, err = validateToken(requestBody.AccessToken)

	if err != nil {
		showResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if models.IsValidImageExtension(requestBody.Image.Filename) {
		uploadFolder := os.Getenv("UPLOAD_FOLDER")
		imagePath := fmt.Sprintf("%s/%s", uploadFolder, requestBody.Image.Filename)

		if err = c.SaveUploadedFile(requestBody.Image, imagePath); err != nil {
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		defer os.Remove(imagePath)
		processedImage, err := models.DecodeImage(imagePath)

		if err != nil {
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		imageData, err := models.CreateImageData(processedImage)

		if err != nil {
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		similarImages, err := env.db.GetSimilarImages(imageData.ImageHash)

		if err != nil {
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"images": similarImages,
		})
	}
}
