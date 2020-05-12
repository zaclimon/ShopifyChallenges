package server

import (
	"UtsuruConcept/controllers"
	"UtsuruConcept/models"
	"log"
)

// Init initializes the server routes.
func Init() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	router := controllers.UtsuruRouter(db)
	router.Run()
}
