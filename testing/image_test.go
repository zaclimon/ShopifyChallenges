package testing

import (
	"UtsuruConcept/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

var imageRows = []string{"id", "file_name", "url", "size", "user_id"}

func TestIsValidImageExtension(t *testing.T) {
	isValidExtensionJpg := models.IsValidImageExtension("test.jpg")
	isValidExtensionJpg2 := models.IsValidImageExtension("test.jpeg")
	isValidExtensionPng := models.IsValidImageExtension("test.png")
	isValidExtensionGif := models.IsValidImageExtension("test.gif")

	if !isValidExtensionJpg || !isValidExtensionJpg2 || !isValidExtensionPng || !isValidExtensionGif {
		t.Error("Valid image extension was not considered valid")
	}

	isValidExtensionRandom := models.IsValidImageExtension("test.txt")
	if isValidExtensionRandom {
		t.Error("Non valid image extension was considered valid.")
	}
}

func TestIsUserImageExists(t *testing.T) {
	sqlmockDb, mock, _ := sqlmock.New()
	defer sqlmockDb.Close()
	gormDb, _ := gorm.Open("mysql", sqlmockDb)
	dbType := &models.DB{DB: gormDb}
	userID := "39e9a9a1-2de0-4a46-a9be-f72da68d0d12"
	imageID := "39e9a9a1-2de0-4a46-a9be-f72da68d0d13"
	url := "https://url.com/test.png"

	t.Run("User has the image", func(t *testing.T) {
		fileName := "test.png"
		rows := sqlmock.NewRows(imageRows).AddRow(imageID, fileName, url, 0, userID)
		mock.ExpectQuery("SELECT").WithArgs(userID, fileName).WillReturnRows(rows)
		result := dbType.IsUserImageExists(userID, fileName)
		if !result {
			t.Error("Can't retrieve an existing image for a user")
		}
	})

	t.Run("User doesn't have the image (Image does not exist)", func(t *testing.T) {
		fileName := "test2.png"
		result := dbType.IsUserImageExists(userID, fileName)
		if result {
			t.Error("Image retrieved for a user when it should not exist")
		}
	})

	t.Run("User doesn't have the image (Image exists)", func(t *testing.T) {
		fileName := "test.png"
		userID2 := "39e9a9a1-2de0-4a46-a9be-f72da68d0d11"
		rows := sqlmock.NewRows(imageRows).AddRow(imageID, fileName, url, 0, userID)
		mock.ExpectQuery("SELECT").WithArgs(userID, fileName).WillReturnRows(rows)
		result := dbType.IsUserImageExists(userID2, fileName)
		if result {
			t.Error("Image retrieved for a user that does not own it")
		}
	})
}
