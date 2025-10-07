package update

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"main.go/db"
	"main.go/middalware"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из токена
	userIDRaw := r.Context().Value(middalware.UserIDKey)
	if userIDRaw == nil {
		http.Error(w, "Токен не содержит пользователя", http.StatusUnauthorized)
		return
	}
	userPhone := userIDRaw.(string)

	// Парсим multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	storeAddress := r.FormValue("store_address")

	// Получаем файл
	file, handler, err := r.FormFile("storeImage")
	var storeImagePath string
	if err == nil {
		defer file.Close()

		dir := "uploads"
		os.MkdirAll(dir, os.ModePerm)

		storeImagePath = filepath.Join(dir, handler.Filename)
		out, err := os.Create(storeImagePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения изображения", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = out.ReadFrom(file)
		if err != nil {
			http.Error(w, "Ошибка записи изображения", http.StatusInternalServerError)
			return
		}
	}

	// Обновляем в базе
	query := `
	UPDATE users
	SET first_name = ?, last_name = ?, store_address = ?, store_image = ?
	WHERE phone = ?
	`

	_, err = db.DB.Exec(query, firstName, lastName, storeAddress, storeImagePath, userPhone)
	if err != nil {
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Профиль обновлён",
	})
}
