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

	db.AutoMigrate(&User{}, &Image{}, &ImageData{})
	user := &User{
		ID:        uuid.New(),
		Email:     "test@test.com",
		Password:  "fdiosafjoias",
		CreatedAt: time.Now(),
	}
	imageData := &ImageData{
		ID:        uuid.New(),
		ImageHash: "ijklofgdsijofsdjio",
	}
	image := &Image{
		ID:       uuid.New(),
		FileName: "Test.jpg",
		Url:      "https://fijaosdfjoids.com/Test.jpg",
		Size:     120583498,
		UserID:   user.ID,
	}

	image.ImageData = *imageData
	user.Images = append(user.Images, *image)
	db.Create(&user)
	fmt.Printf("Is new user? %t", db.NewRecord(user))
	defer db.Close()
}
