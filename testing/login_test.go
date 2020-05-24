package testing

import (
	"UtsuruConcept/controllers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoginEndpoint(t *testing.T) {
	endpoint := "/api/v1/login"

	if err := os.Setenv("JWT_SECRET", "test"); err != nil {
		t.Error("Could not set the secret for JWT tokens")
	}

	t.Run("Valid credentials", func(t *testing.T) {
		dbType, mock, rows := createDbMocks()
		router := controllers.UtsuruRouter(createEnv(dbType))
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
	})

	t.Run("Invalid email input", func(t *testing.T) {
		dbType, mock, _ := createDbMocks()
		router := controllers.UtsuruRouter(createEnv(dbType))
		recorder := httptest.NewRecorder()

		userEmail := "testblah"
		userPassword := "123456"

		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnError(errors.New("Email doesn't exist"))

		loginRequest := controllers.LoginRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(loginRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Error("User have been logged in with an invalid email input")
		}
	})

	t.Run("Unregistered email", func(t *testing.T) {
		dbType, mock, rows := createDbMocks()
		router := controllers.UtsuruRouter(createEnv(dbType))
		recorder := httptest.NewRecorder()

		userEmail := "test2@test.com"
		userPassword := "123456"

		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)

		loginRequest := controllers.LoginRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(loginRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusNotFound {
			t.Error("User with a unregistered email have been logged in")
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		dbType, mock, rows := createDbMocks()
		router := controllers.UtsuruRouter(createEnv(dbType))
		recorder := httptest.NewRecorder()

		userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3b")
		userEmail := "test3@test.com"
		userPassword := "123456"
		userInvalidPassword := "1234567"
		encryptedUserPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

		rows = mock.NewRows([]string{"id", "email", "password", "created_at"})
		rows.AddRow(userID, userEmail, encryptedUserPassword, time.Now())
		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)

		loginRequest := controllers.LoginRequest{
			Email:    userEmail,
			Password: userInvalidPassword,
		}

		requestJson, _ := json.Marshal(loginRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusForbidden {
			t.Error("Could login user with invalid password")
		}
	})

	t.Run("Fields that are not 'username' or 'password", func(t *testing.T) {
		dbType, _, _ := createDbMocks()
		router := controllers.UtsuruRouter(createEnv(dbType))
		recorder := httptest.NewRecorder()

		userEmail := "test@test.com"
		userPassword := "123456"

		randomJson := fmt.Sprintf("{bar: %s, foo: %s}", userEmail, userPassword)

		request := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(randomJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Error("Could execute the request with random fields instead of 'username' and 'password'")
		}
	})

}
