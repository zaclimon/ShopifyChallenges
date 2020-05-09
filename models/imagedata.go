package models

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"image"
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
