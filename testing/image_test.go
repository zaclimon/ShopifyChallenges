package testing

import (
	"UtsuruConcept/models"
	"testing"
)

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
