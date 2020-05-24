package main

import (
	"UtsuruConcept/controllers"
	"UtsuruConcept/models"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// main starts the server and runs the application.
func main() {
	if verifyEnvironmentVariables() {
		Init()
	} else {
		log.Fatal("Cannot run without environment variables!")
	}
}

// verifyEnvironmentVariables checks environment variables has been set, and loads the local one if not.
// It returns true if environment variables could be loaded.
func verifyEnvironmentVariables() bool {
	// The environment variables should already be defined in a container beforehand
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		// Load a .env file in case where we are developing locally
		if err := godotenv.Load(); err != nil {
			log.Fatal("Cannot load environment variables file!")
			return false
		}
	}
	return true
}

// Init initializes the server routes.
func Init() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	env := &controllers.Env{
		Db:   db,
		Mode: controllers.ServerMode,
	}
	router := controllers.UtsuruRouter(env)
	if err = router.Run(); err != nil {
		log.Fatal(err)
	}
}
