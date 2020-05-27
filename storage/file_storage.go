package storage

import (
	"UtsuruConcept/models"
	"cloud.google.com/go/storage"
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Gcs struct {
	*storage.Client
}

type UtsuruFileStorage interface {
	ProcessUpload(user *models.User, files []*multipart.FileHeader, bucketName string, c *gin.Context) ([]string, []string, error)
	CreateUserFolder(folderName string, userID string) error
	UploadToProvider(folderName string, userID string, fileName string) error
	GenerateImageURL(folderName string, userID string, fileName string) (string, error)
}

func InitStorage() (*Gcs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		// The client has not been configured correctly
		return nil, err
	}

	return &Gcs{client}, nil
}
