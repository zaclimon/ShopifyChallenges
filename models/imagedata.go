package models

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"image"
)

type ImageData struct {
	ID        uuid.UUID `gorm:"primary_key"`
	ImageHash uint64
	ImageID   uuid.UUID
}

func (imageData *ImageData) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the imagedata")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

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
