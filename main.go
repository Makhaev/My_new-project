package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"main.go/db"
	"main.go/handlers"
	"main.go/middalware"
	"main.go/models"
	"main.go/profile"
	"main.go/update"
	"main.go/verification"
)

// ddd
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db.Init()

	models.CreateUsers()

	r := chi.NewRouter()

	r.Post("/send_code/", handlers.SendSMS)
	r.Post("/verify_code/", verification.Verification)
	r.Post("/refresh_token/", verification.RefreshToken)
	r.Patch("me/update/", update.UpdateProfile)

	r.Group(func(protected chi.Router) {
		protected.Use(middalware.AuthMidalware)
		protected.Get("/me/", profile.Profile)
		protected.Get("/user/profile", middalware.ProtectedHandler)
	})
	err = http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalf("Ошибка :%v", err)
	}

}
