package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"main.go/db"
	"main.go/handlers"
	"main.go/utils"
)

func main() {
	// 👇 Важно: инициализируем базу данных
	db.Init()

	// Роутер
	r := gin.Default()

	// Роуты
	r.POST("/token/", handlers.Login)
	r.POST("/register/", handlers.Register)

	// Защищённый роут
	protected := r.Group("/user")
	protected.Use(utils.AuthMiddleware())
	protected.GET("/profile/", handlers.GetProfile)

	// Старт сервера
	r.Run(":8080")
}
