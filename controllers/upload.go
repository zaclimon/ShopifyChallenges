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
	"io"
	"net/http"
	"os"
)

type UploadRequest struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}

func Upload(c *gin.Context) {
	var multipartFormRequest UploadRequest
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate that the file is an image.
		// Let's obtain the files first
		form, err := c.MultipartForm()

		if err != nil {
			showResponseError(c, http.StatusBadRequest, err)
			return
		}

		var user models.User

		userID := claims["id"].(string)
		dbObj := db.GetDb()
		dbObj.Where("id = ?", userID).First(&user)
		files := form.File["images"]
		ctx := context.Background()
		client, err := storage.NewClient(ctx)

		if err != nil {
			// The client has not been configured correctly
			showResponseError(c, http.StatusInternalServerError, err)
			return
		}

		bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
		imagesFolderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
		bucket := client.Bucket(bucketName)
		uploadedFiles := make([]string, 0)
		notUploadedFiles := make([]string, 0)

		for _, fileInfo := range files {
			// Verify the extension of the file
			if models.IsValidImageExtension(fileInfo.Filename) && !models.IsUserImageExists(userID, fileInfo.Filename, dbObj) {
				// Upload to cloud storage
				// 1. Ensure that the bucket (folder) for the user exists. Create it if not available
				// 2. Verify if file metadata doesn't exist first.
				// 3. Create the file metadata in the database.
				// 4. Upload the file to the bucket (If filename already exists, don't upload it twice --> Filename based upload instead of ID based)

				userFolder := bucket.Object(fmt.Sprintf("%s/%s/", imagesFolderName, userID))
				_, err = userFolder.Attrs(ctx)

				if err != nil {
					// User folder does not exist yet
					fmt.Println("Creating folder for user")
					folderWriter := userFolder.NewWriter(ctx)
					_, err := folderWriter.Write(make([]byte, 0))
					_ = folderWriter.Close()
					if err != nil {
						// Error happened while trying to create the user folder.
						showResponseError(c, http.StatusInternalServerError, err)
						return
					}
				}
				// Save locally the uploaded file
				fmt.Println("Saving file locally")
				uploadFolder := os.Getenv("UPLOAD_FOLDER")
				destinationPath := fmt.Sprintf("%s/%s", uploadFolder, fileInfo.Filename)
				if err = c.SaveUploadedFile(fileInfo, destinationPath); err != nil {
					showResponseError(c, http.StatusInternalServerError, err)
					return
				}
				// Upload image to GCP
				fmt.Println("Uploading to GCP")
				imageObject := bucket.Object(fmt.Sprintf("%s/%s/%s", imagesFolderName, userID, fileInfo.Filename))
				storageWriter := imageObject.NewWriter(ctx)
				savedFileReader, err := os.Open(destinationPath)
				if err != nil {
					showResponseError(c, http.StatusInternalServerError, err)
					return
				}
				_, err = io.Copy(storageWriter, savedFileReader)
				if err != nil {
					showResponseError(c, http.StatusInternalServerError, err)
					return
				}
				_ = storageWriter.Close()
				_ = savedFileReader.Close()
				// Create database metadata
				// 1. Compute PHash
				decodedImage, err := models.DecodeImage(destinationPath)
				if err != nil {
					fmt.Printf("Affected file: %s\n", fileInfo.Filename)
					showResponseError(c, http.StatusInternalServerError, err)
					_ = savedFileReader.Close()
					_ = os.Remove(destinationPath)
					return
				}
				// We don't want to fill up the storage in the VM/container needlessly so delete the temporary image.
				fmt.Println("Deleting temporary file on system")
				_ = os.Remove(destinationPath)
				// 2. Create ImageData struct
				imageData, err := models.CreateImageData(decodedImage)
				if err != nil {
					showResponseError(c, http.StatusInternalServerError, err)
					return
				}
				// 3. Create Image Struct
				imageModel := models.CreateImage(fileInfo.Filename, fileInfo.Size, *imageData)
				// 4. Link Image to user
				user.Images = append(user.Images, *imageModel)
				uploadedFiles = append(uploadedFiles, fileInfo.Filename)
			} else {
				notUploadedFiles = append(notUploadedFiles, fileInfo.Filename)
			}
		}
		// 5. Update user
		dbObj.Save(&user)
		// 6. Return files that has been uploaded and not uploaded to the user.
		c.JSON(http.StatusOK, gin.H{
			"uploaded_files":     uploadedFiles,
			"not_uploaded_files": notUploadedFiles,
		})
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
