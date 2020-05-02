package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key"`
	Email     string    `gorm:"size:255;unique"`
	Password  string    `gorm:"size:72"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Images    []Image
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the user")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

func InsertNewUser(email string, password string, dbObj *gorm.DB) error {
	if isUserExists(email, dbObj) {
		return errors.New("A user with this email already exists")
	}

	newUser := &User{
		ID:       uuid.New(),
		Email:    email,
		Password: password,
	}

	dbObj.Create(newUser)
	return nil
}

func isUserExists(email string, dbObj *gorm.DB) bool {
	var user User
	dbObj.First(&user, "email = ?", email)
	return user.Email != ""
}
