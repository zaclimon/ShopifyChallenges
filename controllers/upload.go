package controllers

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type UploadRequest struct {
	Token  string                  `form:"token" binding:"required"`
	Images []*multipart.FileHeader `form:"images" binding:"required"`
}

func Upload(c *gin.Context) {
	var requestBody UploadRequest
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(string)
		dbObj := db.GetDb()
		user, bucket, err := initUploadProcess(c, userID, dbObj)
		if err != nil {
			// The body has been defined in initUploadProcess()
			return
		}
		uploadedFiles, notUploadedFiles, err := uploadProcess(user, requestBody.Images, bucket, dbObj, c)

		if err != nil {
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		dbObj.Save(&user)
		c.JSON(http.StatusOK, gin.H{
			"uploaded_files":     uploadedFiles,
			"not_uploaded_files": notUploadedFiles,
		})
	} else {
		showResponseError(c, http.StatusBadRequest, err)
		return
	}
}

func initUploadProcess(c *gin.Context, userID string, dbObj *gorm.DB) (*models.User, *storage.BucketHandle, error) {

	user, err := models.GetUserById(userID, dbObj)

	if err != nil {
		// The user does not exist
		showResponseError(c, http.StatusBadRequest, err)
		return nil, nil, err
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		// The client has not been configured correctly
		showResponseError(c, http.StatusInternalServerError, err)
		return nil, nil, err
	}

	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	bucket := client.Bucket(bucketName)
	return user, bucket, nil
}

func showResponseError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func uploadProcess(user *models.User, files []*multipart.FileHeader, bucket *storage.BucketHandle, dbObj *gorm.DB, c *gin.Context) ([]string, []string, error) {
	imagesFolderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
	uploadedFiles := make([]string, 0)
	notUploadedFiles := make([]string, 0)
	for _, fileInfo := range files {
		userID := user.ID.String()

		if models.IsValidImageExtension(fileInfo.Filename) && !models.IsUserImageExists(userID, fileInfo.Filename, dbObj) {
			userFolder := bucket.Object(fmt.Sprintf("%s/%s/", imagesFolderName, userID))
			ctx := context.Background()
			if _, err := userFolder.Attrs(ctx); err != nil {
				if err = createBucketUserFolder(userFolder); err != nil {
					return nil, nil, err
				}
			}

			if err := saveUploadedFile(fileInfo, c); err != nil {
				return nil, nil, err
			}

			if err := uploadToGCP(bucket, imagesFolderName, userID, fileInfo.Filename); err != nil {
				return nil, nil, err
			}

			imageData, err := generateImageData(fileInfo.Filename)

			if err != nil {
				return nil, nil, err
			}

			imageModel := models.CreateImage(fileInfo.Filename, fileInfo.Size, *imageData)
			user.Images = append(user.Images, *imageModel)
			uploadedFiles = append(uploadedFiles, fileInfo.Filename)
		} else {
			notUploadedFiles = append(notUploadedFiles, fileInfo.Filename)
		}
	}
	return uploadedFiles, notUploadedFiles, nil
}

func createBucketUserFolder(userFolderHandle *storage.ObjectHandle) error {
	folderWriter := userFolderHandle.NewWriter(context.Background())
	_, err := folderWriter.Write(make([]byte, 0))
	_ = folderWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func saveUploadedFile(fileInfo *multipart.FileHeader, c *gin.Context) error {
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	destinationPath := fmt.Sprintf("%s/%s", uploadFolder, fileInfo.Filename)
	if err := c.SaveUploadedFile(fileInfo, destinationPath); err != nil {
		return err
	}
	return nil
}

func uploadToGCP(bucket *storage.BucketHandle, imagesFolderName string, userID string, fileName string) error {
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	imageObject := bucket.Object(fmt.Sprintf("%s/%s/%s", imagesFolderName, userID, fileName))
	storageWriter := imageObject.NewWriter(context.Background())
	savedFileReader, err := os.Open(fmt.Sprintf("%s/%s", uploadFolder, fileName))
	if err != nil {
		return err
	}
	if _, err = io.Copy(storageWriter, savedFileReader); err != nil {
		return err
	}
	defer storageWriter.Close()
	defer savedFileReader.Close()
	return nil
}

func generateImageData(fileName string) (*models.ImageData, error) {
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	destinationPath := fmt.Sprintf("%s/%s", uploadFolder, fileName)
	decodedImage, err := models.DecodeImage(destinationPath)
	defer os.Remove(destinationPath)

	if err != nil {
		return nil, err
	}

	imageData, err := models.CreateImageData(decodedImage)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
