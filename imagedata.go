package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ImageData struct {
	ID        uuid.UUID `gorm:"primary_key"`
	ImageHash string
	ImageId   uuid.UUID
}

func (imageData *ImageData) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the imagedata")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}
