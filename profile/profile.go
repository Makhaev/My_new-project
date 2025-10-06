package profile

import (
	"encoding/json"
	"net/http"

	"main.go/middalware"
	"main.go/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middalware.UserIDKey)
	if userIDRaw == nil {
		http.Error(w, "Токен не содержит пользователя", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		http.Error(w, "Неверный формат пользователя в токене", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByPhone(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
