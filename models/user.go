package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:72;not null"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Images    []Image
}

var DbUserNotFoundError = errors.New("The user does not exist")
var DbDuplicatedEmailError = errors.New("A user with this email already exists")
var InvalidCredentialsError = errors.New("The credentials entered are not valid")
var PasswordHashError = errors.New("An error happened while hashing the password")

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	generatedUuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Errorf("error while generating UUID for the user")
		return err
	}

	return scope.SetColumn("ID", generatedUuid)
}

func CreateAndInsertNewUser(email string, password string, dbObj *gorm.DB) (*User, error) {
	if isUserExists(email, dbObj) {
		return nil, DbDuplicatedEmailError
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedPassword,
	}

	dbObj.Create(&newUser)
	return newUser, nil
}

func GetUserByEmail(email string, password string, dbObj *gorm.DB) (*User, error) {
	if !isUserExists(email, dbObj) {
		return nil, DbUserNotFoundError
	}
	var user User
	dbObj.First(&user, "email = ?", email)

	if validatePassword(user.Password, password) {
		return &user, nil
	}

	return nil, InvalidCredentialsError
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", PasswordHashError
	}
	return string(bytes), nil
}

func isUserExists(email string, dbObj *gorm.DB) bool {
	var user User
	dbObj.First(&user, "email = ?", email)
	return user.Email != ""
}

func validatePassword(passwordHash string, passwordCandidate string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordCandidate)) == nil
}

func UpdateUser(dbObj *gorm.DB, user User) {
	dbObj.Update(user)
}
