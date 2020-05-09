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

// Error when a user does not exist in the database
var DbUserNotFoundError = errors.New("The user does not exist")

// Error when a user with a given email already exists
var DbDuplicatedEmailError = errors.New("A user with this email already exists")

// Error when credentials are not valid
var InvalidCredentialsError = errors.New("The credentials entered are not valid")

// Error when the password could not be hashed
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

// CreateAndInsertNewUser creates a user type and inserts it on the database based on an email, password and database object.
func CreateAndInsertNewUser(email string, password string, dbObj *gorm.DB) (*User, error) {
	if isUserEmailExists(email, dbObj) {
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

// GetUserByEmail retrieves a user from the database based on its email, password, database object.
func GetUserByEmail(email string, password string, dbObj *gorm.DB) (*User, error) {
	if !isUserEmailExists(email, dbObj) {
		return nil, DbUserNotFoundError
	}
	var user User
	dbObj.First(&user, "email = ?", email)

	if validatePassword(user.Password, password) {
		return &user, nil
	}

	return nil, InvalidCredentialsError
}

// GetUserById retrieves a user from the database based on its id and database object.
func GetUserById(id string, dbObj *gorm.DB) (*User, error) {
	if !isUserIdExists(id, dbObj) {
		return nil, DbUserNotFoundError
	}
	var user User
	dbObj.First(&user, "id = ?", id)
	return &user, nil
}

// hashPassword hashes the password of the user
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", PasswordHashError
	}
	return string(bytes), nil
}

// isUserEmailExists verifies a user with the given email exists.
func isUserEmailExists(email string, dbObj *gorm.DB) bool {
	var user User
	dbObj.First(&user, "email = ?", email)
	return user.Email != ""
}

func isUserIdExists(id string, dbObj *gorm.DB) bool {
	var user User
	dbObj.First(&user, "id = ?", id)
	return user.ID != uuid.Nil
}

func validatePassword(passwordHash string, passwordCandidate string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordCandidate)) == nil
}
