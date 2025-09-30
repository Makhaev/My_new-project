package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main.go/models"
)

// шаг
// Иницилизация базы данных
// создается соеденение Sqlite

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gopher.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
}
