package souko

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"souko/models"
)

func ConfigureDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		panic("Could not migrate database objects")
	}
	models.ConfigureProductDao(db)
	return db
}
