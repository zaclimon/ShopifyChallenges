package main

import (
	"UtsuruConcept/db"
	"UtsuruConcept/server"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if verifyEnvironmentVariables() {
		db.Init()
		server.Init()
		db.Stop()
	} else {
		log.Fatal("Cannot run without environment variables!")
	}
}

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
