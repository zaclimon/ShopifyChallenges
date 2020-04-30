package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id        uuid.UUID `gorm:"primary_key"`
	Email     string    `gorm:"size:255;unique"`
	Password  string    `gorm:"size:72"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the user")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}
