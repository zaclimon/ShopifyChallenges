package db

import (
	"UtsuruConcept/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbObj *gorm.DB

func Init() {
	if dbObj == nil {
		db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/mysql?&parseTime=True&loc=UTC")

		if err != nil {
			fmt.Println("Error while trying to initialize the database", err)
		}

		dbObj = db
		dbObj.AutoMigrate(&models.User{}, &models.Image{}, &models.ImageData{})
	} else {
		fmt.Println("Database has already been initialized")
	}
}

func GetDb() *gorm.DB {
	if dbObj != nil {
		return dbObj
	}
	fmt.Println("Please initialize database object")
	return nil
}

func Stop() {
	if dbObj != nil {
		defer dbObj.Close()
	}
}
