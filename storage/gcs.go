package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
)

type Gcs struct {
	*storage.Client
}

type UtsuruFileStorage interface {
	// IsUserFolderExists verifies if the file storage folder for a user exists
	IsUserFolderExists(userID string) bool
	// CreateUserFolder creates a new folder for a user inside a file storage bucket.
	CreateUserFolder(userID string) error
	// UploadToProvider imports the uploaded file for the user
	UploadToProvider(userID string, fileName string) error
	// GenerateImageURL retrieves the URL that will be used for public sharing on the cloud storage service
	GenerateImageURL(userID string, fileName string) (string, error)
}

// InitStorage initializes the objects for Utsuru
func InitStorage() (*Gcs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		// The client has not been configured correctly
		return nil, err
	}

	return &Gcs{client}, nil
}

func (gcs *Gcs) IsUserFolderExists(userID string) bool {
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	folderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
	bucket := gcs.Bucket(bucketName)
	userFolder := bucket.Object(fmt.Sprintf("%s/%s/", folderName, userID))
	ctx := context.Background()
	_, err := userFolder.Attrs(ctx)
	return err == nil
}

func (gcs *Gcs) CreateUserFolder(userID string) error {
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	folderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
	bucket := gcs.Bucket(bucketName)
	folderHandle := bucket.Object(fmt.Sprintf("%s/%s", folderName, userID))
	folderWriter := folderHandle.NewWriter(context.Background())
	_, err := folderWriter.Write(make([]byte, 0))
	_ = folderWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (gcs *Gcs) UploadToProvider(userID string, fileName string) error {
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	imagesFolderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
	uploadFolder := os.Getenv("UPLOAD_FOLDER")
	bucket := gcs.Bucket(bucketName)
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

func (gcs *Gcs) GenerateImageURL(userID string, fileName string) (string, error) {
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	imagesFolderName := os.Getenv("CLOUD_STORAGE_IMAGES_FOLDER")
	bucket := gcs.Bucket(bucketName)
	imageObject := bucket.Object(fmt.Sprintf("%s/%s/%s", imagesFolderName, userID, fileName))
	imageAcl := imageObject.ACL()

	// Set image as public
	ctx := context.Background()
	if err := imageAcl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	imageAttrs, err := imageObject.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return imageAttrs.MediaLink, nil
}
