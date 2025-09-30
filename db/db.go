package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main.go/models"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gopher.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}
