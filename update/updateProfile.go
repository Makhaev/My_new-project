package update

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	var firstName, lastName, storeAddress, storeImagePath string

	contentType := r.Header.Get("Content-Type")

	// 🧩 1. Если JSON
	if strings.HasPrefix(contentType, "application/json") {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
			return
		}
		firstName = data["first_name"]
		lastName = data["last_name"]
		storeAddress = data["store_address"]
	}

	// 🧩 2. Если multipart/form-data (форма + изображение)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			return
		}

		firstName = r.FormValue("first_name")
		lastName = r.FormValue("last_name")
		storeAddress = r.FormValue("store_address")

		file, handler, err := r.FormFile("storeImage")
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
	}

	// 🧩 3. Обновляем в базе
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, store_address = $3, store_image = $4
	WHERE phone = $5
	`

	_, err := db.DB.Exec(query, firstName, lastName, storeAddress, storeImagePath, userPhone)
	if err != nil {
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}

	// 🧩 4. Ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Профиль обновлён",
	})
}
