package verification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"main.go/db"
	"main.go/generathionToken"
)

type VerifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type VerifyResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
	Token  string `json:"token,omitempty"`
}

func Verification(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Неверный формат JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var createdAt time.Time

	query := `SELECT created_at FROM sms_codes WHERE phone = ? AND code = ? ORDER BY created_at DESC LIMIT 1`
	err := db.DB.QueryRow(query, req.Phone, req.Code).Scan(&createdAt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Error: "Неверный код"})
		return
	}

	if time.Since(createdAt) > 5*time.Minute {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Error: "Код истёк"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM sms_codes WHERE phone = ? AND code = ?", req.Phone, req.Code)
	if err != nil {
		http.Error(w, `{"error":"Ошибка удаления кода"}`, http.StatusInternalServerError)
		return
	}

	token, err := generathionToken.GenerateToken(req.Phone)
	if err != nil {
		http.Error(w, `{"error":"Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}

	refreshToken, err := generathionToken.GenerateRefreshToken(req.Phone)
	if err != nil {
		http.Error(w, `{"error":"Ошибка генерации refresh токена"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Access:", token)
	fmt.Println("✅ Refresh:", refreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "verified",
		"access":  token,
		"refresh": refreshToken,
	})
}
