package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is a type that identifies a user for Utsuru.
type User struct {
	ID        uuid.UUID `gorm:"primary_key"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:72;not null"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Images    []Image
}

// Error when a user does not exist in the database.
var DbUserNotFoundError = errors.New("The user does not exist")

// Error when a user with a given email already exists.
var DbDuplicatedEmailError = errors.New("A user with this email already exists")

// Error when credentials are not valid.
var InvalidCredentialsError = errors.New("The credentials entered are not valid")

// Error when the password could not be hashed.
var PasswordHashError = errors.New("An error happened while hashing the password")

// BeforeCreate is a function called by Gorm for preliminary processing before inserting a new object in the database.
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the user")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

func (db *DB) InsertOrUpdateUser(user *User) error {
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewUser(email string, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	newUser := &User{
		Email:    email,
		Password: hashedPassword,
	}
	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	var user User
	db.First(&user, "email = ?", email)
	if user.Email == "" {
		return nil, DbUserNotFoundError
	}
	return &user, nil
}

func (db *DB) GetUserById(userID string) (*User, error) {
	var user User
	db.First(&user, "id = ?", userID)
	if user.ID == uuid.Nil {
		return nil, DbUserNotFoundError
	}
	return &user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", PasswordHashError
	}
	return string(bytes), nil
}

func (db *DB) IsUserEmailExists(email string) bool {
	if _, err := db.GetUserByEmail(email); err != nil {
		return false
	}
	return true
}

func (db *DB) IsUserIdExists(id string) bool {
	if _, err := db.GetUserById(id); err != nil {
		return false
	}
	return true
}

func ValidatePassword(passwordHash string, passwordCandidate string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordCandidate))
}
