package models

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"image"
	"os"
	"strconv"
)

// ImageData is a type for registering specific image metadata.
type ImageData struct {
	ID        uuid.UUID `gorm:"primary_key"`
	ImageHash uint64
	ImageID   uuid.UUID
}

// BeforeCreate is a function called by Gorm for preliminary processing before inserting a new object in the database.
func (imageData *ImageData) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the imagedata")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

// CreateImage makes an ImageData type based on a processed image.
// It returns an ImageData type or an error if metadata could not be processed
func CreateImageData(image image.Image) (*ImageData, error) {
	hash, err := goimagehash.PerceptionHash(image)
	if err != nil {
		return nil, err
	}

	newImageData := &ImageData{
		ImageHash: hash.GetHash(),
	}

	return newImageData, nil
}

func (db *DB) GetSimilarImages(imageHash uint64) (*[]Image, error) {
	var images []Image
	hashThreshold, err := strconv.Atoi(os.Getenv("PHASH_THRESHOLD"))

	if err != nil {
		return nil, err
	}

	// Not sure if a raw SQL query is the right way to go but it works...
	sqlQuery := fmt.Sprintf("SELECT images.* from images, image_data WHERE images.id = image_data.image_id AND bit_count(%d ^ image_data.image_hash) <= %d", imageHash, hashThreshold)
	db.Raw(sqlQuery).Find(&images)
	return &images, nil
}

// generateImageData creates image metadata that can be used for further processing.
// It returns an error if the image could not be decoded or if it could not extract metadata from the image.
func GenerateImageData(filePath string) (*ImageData, error) {
	decodedImage, err := DecodeImage(filePath)

	if err != nil {
		return nil, err
	}

	imageData, err := CreateImageData(decodedImage)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
