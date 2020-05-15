package testing

import (
	"UtsuruConcept/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func initMocks() (*models.DB, sqlmock.Sqlmock, *sqlmock.Rows) {
	sqlmockDb, mock, _ := sqlmock.New()
	gormDb, _ := gorm.Open("mysql", sqlmockDb)
	dbType := &models.DB{DB: gormDb}
	rows := mock.NewRows([]string{"id", "email", "password", "created_at"})
	return dbType, mock, rows
}

func TestGetUserByEmail(t *testing.T) {
	dbType, mock, rows := initMocks()
	userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
	userEmail := "test@test.com"
	userPassword := "abc123"
	rows.AddRow(userID, userEmail, userPassword, time.Now())

	t.Run("User exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)
		user, err := dbType.GetUserByEmail(userEmail)
		if err != nil {
			t.Error("Could not retrieve valid user correctly with mock")
		}

		if user.ID != userID || user.Email != userEmail || user.Password != userPassword {
			t.Error("Mismatch between mocked user and retrieved user")
		}
	})

	t.Run("User does not exist", func(t *testing.T) {
		email := "test2@test.com"
		_, err := dbType.GetUserByEmail(email)
		if err == nil {
			t.Error("Could retrieve a user based on a unregistered email.")
		}
	})
}

func TestIsUserEmailExists(t *testing.T) {
	dbType, mock, rows := initMocks()
	userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
	userEmail := "test@test.com"
	userPassword := "abc123"
	rows.AddRow(userID, userEmail, userPassword, time.Now())

	t.Run("User exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(userEmail).WillReturnRows(rows)
		if !dbType.IsUserEmailExists(userEmail) {
			t.Error("Could not retrieve valid user correctly with mock")
		}
	})

	t.Run("User does not exist", func(t *testing.T) {
		email := "test2@test.com"
		if dbType.IsUserEmailExists(email) {
			t.Error("Could retrieve a user based on a unregistered email.")
		}
	})
}

func TestCreateNewUser(t *testing.T) {
	userEmail := "test@test.com"
	userPassword := "123456"
	user, err := models.CreateNewUser(userEmail, userPassword)

	if err != nil {
		t.Error("Could not create a generic user")
	}

	if user.Password != userPassword {
		// Password is probably hashed
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword)); err != nil {
			t.Error("Hashed password and user password mismatch")
		}
	}
}

func TestGetUserById(t *testing.T) {
	dbType, mock, rows := initMocks()
	userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
	userEmail := "test@test.com"
	userPassword := "abc123"
	rows.AddRow(userID, userEmail, userPassword, time.Now())

	t.Run("User exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rows)
		user, err := dbType.GetUserById(userID.String())
		if err != nil {
			t.Error("Could not retrieve valid user correctly with mock")
		}

		if user.ID != userID || user.Email != userEmail || user.Password != userPassword {
			t.Error("Mismatch between mocked user and retrieved user")
		}
	})

	t.Run("User does not exist", func(t *testing.T) {
		id, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3b")
		_, err := dbType.GetUserById(id.String())
		if err == nil {
			t.Error("Could retrieve a user based on its ID")
		}
	})
}

func TestIsUserIdExists(t *testing.T) {
	dbType, mock, rows := initMocks()
	userID, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3a")
	userEmail := "test@test.com"
	userPassword := "abc123"
	rows.AddRow(userID, userEmail, userPassword, time.Now())

	t.Run("User exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(userID).WillReturnRows(rows)
		if !dbType.IsUserIdExists(userID.String()) {
			t.Error("Could not retrieve valid user correctly with mock")
		}
	})

	t.Run("User does not exist", func(t *testing.T) {
		id, _ := uuid.Parse("47ff6ccc-5bed-4e9a-8b0f-3ea995c7ac3b")
		if dbType.IsUserIdExists(id.String()) {
			t.Error("Could retrieve a user based on a unregistered email.")
		}
	})
}

func TestValidatePassword(t *testing.T) {
	password := "test1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Error("Could not hash the password")
	}

	if err = models.ValidatePassword(string(hashedPassword), password); err != nil {
		t.Error("Could not validate the password")
	}
}
