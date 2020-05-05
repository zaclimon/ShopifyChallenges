package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
)

type Image struct {
	ID        uuid.UUID `gorm:"primary_key"`
	FileName  string    `gorm:"not null"`
	Url       string
	Size      uint64
	UserID    uuid.UUID
	ImageData ImageData
}

func (image *Image) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the image")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

func IsValidImageExtension(fileName string) bool {
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) >= 1 {
		fileExtension := fileParts[len(fileParts)-1]
		return fileExtension == "jpg" || fileExtension == "jpeg" || fileExtension == "png" || fileExtension == "gif"
	}
	return false
}

func IsUserImageExists(userID string, fileName string, dbObj *gorm.DB) bool {
	var image Image
	dbObj.First(&image, "user_id = ? AND file_name = ?", userID, fileName)
	return image.FileName != ""
}
