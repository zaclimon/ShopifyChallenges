package testing

import (
	"UtsuruConcept/controllers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRegisterEndpoint(t *testing.T) {
	dbType, mock, rows := createDbMocks()
	router := controllers.UtsuruRouter(createEnv(dbType))
	endpoint := "/api/v1/register"

	t.Run("Valid credentials", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userEmail := "test@test.com"
		userPassword := "123456"

		mock.ExpectBegin()
		mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		registerRequest := controllers.RegisterRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(registerRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Error("Could not register user while using valid credentials")
		}
	})

	t.Run("Invalid email input", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userEmail := "testblah"
		userPassword := "123456"

		registerRequest := controllers.RegisterRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(registerRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Error("User have been created with an invalid email input")
		}
	})

	t.Run("Password less than 6 characters", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userEmail := "test@test.com"
		userPassword := "1234"

		registerRequest := controllers.RegisterRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(registerRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Error("User have been created with a password smaller than 6 characters.")
		}
	})

	t.Run("Fields that are not 'username' or 'password", func(t *testing.T) {
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

	t.Run("Duplicate email", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
		userEmail := "test@test.com"
		userPassword := "123456"

		rows.AddRow(userID, userEmail, userPassword, time.Now())
		mock.ExpectBegin()
		mock.ExpectExec("INSERT").WillReturnError(errors.New("User with email already exists (expected)"))

		registerRequest := controllers.RegisterRequest{
			Email:    userEmail,
			Password: userPassword,
		}

		requestJson, _ := json.Marshal(registerRequest)
		request := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestJson))
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusConflict {
			t.Error("Two users with the same email have been created")
		}
	})
}
