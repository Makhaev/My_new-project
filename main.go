package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"main.go/db"
	"main.go/handlers"
	"main.go/middalware"
	"main.go/models"
	"main.go/profile"
	"main.go/update"
	"main.go/verification"
)

func main() {

	db.Init()

	models.CreateUsers()

	r := chi.NewRouter()

	r.Post("/send_code/", handlers.SendSMS)
	r.Post("/verify_code/", verification.Verification)
	r.Post("/refresh_token/", verification.RefreshToken)
	r.Get("/categories/", handlers.GetCategoriesHandler)
	r.Get("/categories/{id}/", handlers.GetCategoryHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(middalware.AuthMidalware)
		protected.Get("/me/", profile.Profile)
		protected.Patch("/me/update/", update.UpdateProfile) // ← исправлено
		protected.Get("/user/profile/", middalware.ProtectedHandler)
		protected.Post("/categories/", handlers.CreateCategoryHandler)
		protected.Patch("/categories/{id}/", handlers.UpdateCategoryHandler)
		protected.Delete("/categories/{id}/", handlers.DeleteCategoryHandler)
	})

	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalf("Ошибка :%v", err)
	}

}
