package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"main.go/db"
	"main.go/generathion"
)

type SMSRequest struct {
	Phone  string `json:"phone"`
	Text   string `json:"text"`
	Sender string `json:"sender"`
}

func SendSMS(w http.ResponseWriter, r *http.Request) {
	var req SMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	code := generathion.GenarathionCode(6)
	expiresAt := time.Now().Add(5 * time.Minute)

	query := `INSERT INTO sms_codes(phone,code,expires_at) VALUES(?,?)`

	_, err := db.DB.Exec(query, req.Phone, code, expiresAt)
	if err != nil {
		http.Error(w, "Ошибка сохранение в базу данных", http.StatusBadRequest)
	}

	// apiKey := "WGZT693J3OWKRM8W47DPS275631Y45SS3V4TPA9J88U1QAS7572C41F53QIEO4OK"

	apiKey := os.Getenv("SMS_API_KEY")
	if apiKey == "" {
		http.Error(w, "API ключ не найден", http.StatusInternalServerError)
		return

	}

	req.Text = fmt.Sprintln("Ваш код авторизации ", code)
	fmt.Println("код отправлен ", req.Text)
	if req.Phone == "+79659628225" {
		encodedText := url.QueryEscape(req.Text)
		smsURL := fmt.Sprintf(
			"http://smspilot.ru/api.php?send=%s&to=%s&from=%s&apikey=%s&format=json",
			encodedText, req.Phone, req.Sender, apiKey,
		)

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(smsURL)
		if err != nil {
			http.Error(w, "ошибка при отправке запроса", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения ответа от SMSPilot", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}

}
