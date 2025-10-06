package profile

import (
	"encoding/json"
	"net/http"

	"main.go/middalware"
	"main.go/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	phoneRaw := r.Context().Value(middalware.UserIDKey)
	if phoneRaw == nil {
		http.Error(w, "Токен не содержит пользователя", http.StatusUnauthorized)
		return
	}

	phone := phoneRaw.(string)

	user, err := models.GetUserByPhone(phone)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
