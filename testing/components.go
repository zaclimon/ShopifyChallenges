package testing

import (
	"UtsuruConcept/controllers"
	"UtsuruConcept/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

// createDbMocks generates appropriate mocks for interacting with the database
func createDbMocks() (*models.DB, sqlmock.Sqlmock, *sqlmock.Rows) {
	sqlmockDb, mock, _ := sqlmock.New()
	gormDb, _ := gorm.Open("mysql", sqlmockDb)
	dbType := &models.DB{DB: gormDb}
	rows := mock.NewRows([]string{"id", "email", "password", "created_at"})
	return dbType, mock, rows
}

// createEnv generates the environment files required for an environment
func createEnv(db *models.DB) *controllers.Env {
	return &controllers.Env{
		Db:   db,
		Mode: controllers.TestMode,
	}
}
