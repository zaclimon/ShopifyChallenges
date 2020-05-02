package main

import (
	"UtsuruConcept/db"
	"UtsuruConcept/models"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db.Init()
	dbObj := db.GetDb()
	err := models.InsertNewUser("test2", "test123", dbObj)

	if err != nil {
		fmt.Println("Something happened while inserting new user")
	}
	db.Stop()
}
