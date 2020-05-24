package testing

import (
	"UtsuruConcept/controllers"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLoginEndpoint(t *testing.T) {
	dbType, mock, rows := initMocks()
	router := controllers.UtsuruRouter(dbType)
	endpoint := "/api/v1/login"

	if err := os.Setenv("JWT_SECRET", "test"); err != nil {
		t.Error("Could not set the secret for JWT tokens")
	}

	t.Run("Valid credentials", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
		userEmail := "test@test.com"
		userPassword := "123456"
		encryptedUserPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

		rows.AddRow(userID, userEmail, encryptedUserPassword, time.Now())
		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)

		loginRequest := controllers.LoginRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(loginRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Error("Could not register user while using valid credentials")
		}
	})
}
