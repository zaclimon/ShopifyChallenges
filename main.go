package main

import (
	"UtsuruConcept/db"
	"UtsuruConcept/server"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db.Init()
	server.Init()
	db.Stop()
}
