package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

// Interface that abstracts database related operations for UtsuruConcept.
type UtsuruDataStore interface {
	// GetSimilarImages retrieves images that are "similar" based on a given image hash.
	GetSimilarImages(imageHash uint64) (*[]Image, error)
	// GetUserByEmail retrieves a user from the database based on its email.
	GetUserByEmail(email string) (*User, error)
	// GetUserById retrieves a user from the database based on its id.
	GetUserById(userID string) (*User, error)
	// InsertsOrUpdateUser inserts or updates the given user information into the database.
	InsertOrUpdateUser(user *User) error
	// IsUserImageExists validates whether a user has a given image.
	IsUserImageExists(userID string, fileName string) bool
	// IsUserEmailExists verifies a user with the given email exists.
	IsUserEmailExists(email string) bool
	// isUserIdExists verifies if a user with the given id exists.
	IsUserIdExists(userID string) bool
}

type DB struct {
	*gorm.DB
}

func InitDB() (*DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&parseTime=True&loc=UTC", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Image{}, &ImageData{})
	return &DB{db}, nil
}
