package verification

import (
	"encoding/json"
	"net/http"
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
	Token  string `json:"token"`
}

func Verification(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Неверный формат JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var createdAt time.Time

	// Получаем время создания кода
	query := `SELECT created_at FROM sms_codes WHERE phone = ? AND code = ? ORDER BY created_at DESC LIMIT 1`
	err := db.DB.QueryRow(query, req.Phone, req.Code).Scan(&createdAt)
	if err != nil {
		// код не найден
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Error: "Неверный код"})
		return
	}

	// Проверка: не истёк ли срок действия (больше 5 минут)
	if time.Since(createdAt) > 5*time.Minute {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Error: "Код истёк"})
		return
	}

	// Удаляем использованный код
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

	user := models.UserStoreInfo{
		StorePhone: req.Phone,
	}

	user.CreateUser()
	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerifyResponse{
		Status: "verified",
		Token:  token,
	})
}
