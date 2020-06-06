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

// UploadRequest handles upload information like the images to upload.
type UploadRequest struct {
	AccessToken string                  `form:"access_token" binding:"required"`
	Images      []*multipart.FileHeader `form:"images" binding:"required"`
}

// Upload uploads one or more images to Google Cloud and creates associated metadata entries in the database when using
// the "/upload" endpoint.
func (env *Env) Upload(c *gin.Context) {
	var requestBody UploadRequest
	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	userID, err := validateToken(requestBody.AccessToken)
	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}
	user, err := env.Db.GetUserById(userID)

	if err != nil {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}

	uploadedFiles, notUploadedFiles, err := processUpload(user, requestBody.Images, env, c)
	if err != nil {
		showResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if err = env.Db.InsertOrUpdateUser(user); err != nil {
		showResponseError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uploaded_files":     uploadedFiles,
		"not_uploaded_files": notUploadedFiles,
	})
}

// processUpload handles the main processing of uploading files to Google Cloud Storage.
// It returns an array of uploaded pictures, pictures that has not been uploaded and an error if something happened.
//
// Note: Two pictures cannot have the same filename in the repository. As such if one is tried to be uploaded, it will
// not be uploaded.
func processUpload(user *models.User, files []*multipart.FileHeader, env *Env, c *gin.Context) ([]string, []string, error) {
	uploadedFiles := make([]string, 0)
	notUploadedFiles := make([]string, 0)
	for _, fileInfo := range files {
		userID := user.ID.String()

		if models.IsValidImageExtension(fileInfo.Filename) && !env.Db.IsUserImageExists(userID, fileInfo.Filename) {
			if !env.Gcs.IsUserFolderExists(userID) {
				if err := env.Gcs.CreateUserFolder(userID); err != nil {
					return nil, nil, err
				}
			}

			fmt.Println("Saving file locally")
			if err := saveUploadedFile(fileInfo, c); err != nil {
				return nil, nil, err
			}

			fmt.Println("Saving file to provider")
			if err := env.Gcs.UploadToProvider(userID, fileInfo.Filename); err != nil {
				return nil, nil, err
			}

			fmt.Println("Generating image URL")
			imageUrl, err := env.Gcs.GenerateImageURL(userID, fileInfo.Filename)

			if err != nil {
				return nil, nil, err
			}

			imagePath := fmt.Sprintf("%s/%s", os.TempDir(), fileInfo.Filename)
			imageData, err := models.GenerateImageData(imagePath)

			fmt.Println("Deleting local image")
			_ = os.Remove(imagePath)

			if err != nil {
				return nil, nil, err
			}

			imageModel := models.CreateImage(fileInfo.Filename, fileInfo.Size, *imageData)
			imageModel.Url = imageUrl
			user.Images = append(user.Images, *imageModel)
			uploadedFiles = append(uploadedFiles, fileInfo.Filename)
		} else {
			notUploadedFiles = append(notUploadedFiles, fileInfo.Filename)
		}
	}
	return uploadedFiles, notUploadedFiles, nil
}

// showResponseError writes to the response body the error retrieved when executing functions
func showResponseError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
