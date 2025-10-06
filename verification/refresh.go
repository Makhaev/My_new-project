package verification

import (
	"encoding/json"
	"net/http"

	"main.go/generathionToken"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Error       string `json:"error,omitempty"`
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Неверный формат JSON"}`, http.StatusBadRequest)
		return
	}

	phone, err := generathionToken.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, `{"error":"Неверный refresh токен"}`, http.StatusUnauthorized)
		return
	}

	newAccess, err := generathionToken.GenerateToken(phone)
	if err != nil {
		http.Error(w, `{"error":"Ошибка генерации access токена"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(RefreshResponse{
		AccessToken: newAccess,
	})
}
