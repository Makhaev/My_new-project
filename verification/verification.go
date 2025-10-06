package verification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"main.go/db"
	"main.go/generathionToken"
	"main.go/models"
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
		fmt.Println("Ошибка генерации токена:", err)
		http.Error(w, `{"error":"Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Токен успешно создан:", token)

	user := models.UserStoreInfo{StorePhone: req.Phone}
	if err := user.CreateUser(); err != nil && !strings.Contains(err.Error(), "UNIQUE") {
		fmt.Println("Ошибка при создании пользователя:", err)
		http.Error(w, `{"error":"Ошибка создания пользователя"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(VerifyResponse{
		Status: "verified",
		Token:  token,
	})
}
