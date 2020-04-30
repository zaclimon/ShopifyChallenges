package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func main() {
	db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/mysql?&parseTime=True&loc=UTC")

	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&User{})
	user := &User{
		Id:        uuid.New(),
		Email:     "test@test.com",
		Password:  "fdiosafjoias",
		CreatedAt: time.Now(),
	}
	db.NewRecord(user)
	db.Create(&user)
	defer db.Close()
}
