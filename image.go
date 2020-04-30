package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Image struct {
	Id       uuid.UUID `gorm:"primary_key"`
	FileName string    `gorm:"not null"`
	Url      string
	Size     uint64
	UserId   uuid.UUID
}

func (image *Image) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the image")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}