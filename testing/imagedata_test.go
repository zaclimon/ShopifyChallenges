package testing

import (
	"UtsuruConcept/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestCreateImageData(t *testing.T) {
	t.Run("Valid image", func(t *testing.T) {
		width := 100
		height := 100
		upLeft := image.Point{X: 0, Y: 0}
		downRight := image.Point{X: width, Y: height}
		img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: downRight})

		// Set a single blue point at the middle of the image
		img.Set(width/2, height/2, color.RGBA{A: 255, R: 0, G: 0, B: 255})
		imageData, err := models.CreateImageData(img)

		if err != nil {
			t.Error("Error while creating image data for a generated image")
		}

		if imageData.ImageHash != 11053635145317115545 {
			t.Error("Image hash of generated image is not valid")
		}
	})

	t.Run("Nil image", func(t *testing.T) {
		_, err := models.CreateImageData(nil)
		if err == nil {
			t.Error("Image metadata created when the given image is nil")
		}
	})
}

func TestGetSimilarImages(t *testing.T) {
	// A better test would probably be to validate the BIT_COUNT function used for computing the hamming distance
	// between the images. However, we would need integration tests for that so let's mock our queries for the
	// time being.
	imageRows := []string{"id", "file_name", "url", "size", "user_id"}
	sqlmockDb, mock, _ := sqlmock.New()
	defer sqlmockDb.Close()
	gormDb, _ := gorm.Open("mysql", sqlmockDb)
	dbType := &models.DB{DB: gormDb}

	imageID := "39e9a9a1-2de0-4a46-a9be-f72da68d0d13"
	fileName := "test.png"
	url := "http://test.com/test.png"
	size := 0
	userID := "39e9a9a1-2de0-4a46-a9be-f72da68d0d11"

	rows := sqlmock.NewRows(imageRows).AddRow(imageID, fileName, url, size, userID)

	t.Run("Without threshold set", func(t *testing.T) {
		_, err := dbType.GetSimilarImages(1234567890)
		if err == nil {
			t.Error("Similar images could be obtained without setting a threshold.")
		}
	})

	// We would be better to use dependency injection rather than setting manually environment variables...
	if err := os.Setenv("PHASH_THRESHOLD", "0"); err != nil {
		t.Error("Error while setting the temporary hash threshold")
	}

	t.Run("With threshold, but no similar images", func(t *testing.T) {
		images, _ := dbType.GetSimilarImages(000000)
		if len(*images) > 0 {
			t.Error("Getting similar images when that should not be the case")
		}
	})

	// Since a raw SQL query has been made, we need to mock the rows before executing it
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	t.Run("With threshold, but with similar images", func(t *testing.T) {
		images, _ := dbType.GetSimilarImages(0123456)
		if len(*images) == 0 {
			t.Error("Not getting any similar images when that should be the case")
		}
	})
}

func TestGenerateImageData(t *testing.T) {
	dir, err := os.Getwd()

	if err != nil {
		t.Error("Could not get testing directory")
	}

	imagePath := fmt.Sprintf("%s/files/test_picture.jpg", dir)
	textPath := fmt.Sprintf("%s/files/test_text.txt", dir)

	t.Run("Generate ImageData with valid image", func(t *testing.T) {
		imageData, err := models.GenerateImageData(imagePath)
		if err != nil {
			t.Error("Could not get a real image metadata")
		}

		if imageData.ImageHash == 000000000 {
			t.Error("Invalid image hash for a valid image")
		}
	})

	t.Run("Generate ImageData with invalid file (not an image)", func(t *testing.T) {
		_, err := models.GenerateImageData(textPath)
		if err == nil {
			t.Error("ImageData generated while file is not an image")
		}
	})

	t.Run("Generate ImageData with no path specified", func(t *testing.T) {
		_, err := models.GenerateImageData("")
		if err == nil {
			t.Error("ImageData generated while no path has been specified")
		}
	})

	t.Run("Generate ImageData with invalid path specified", func(t *testing.T) {
		invalidPath := fmt.Sprintf("%s/something/is/wrong", dir)
		_, err := models.GenerateImageData(invalidPath)
		if err == nil {
			t.Error("ImageData generated while an invalid path has been specified")
		}
	})
}
