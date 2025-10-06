package middalware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"main.go/generathionToken"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMidalware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Токен не предоставлен", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := generathionToken.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey)
	fmt.Println(userID)

	if userID == nil {
		http.Error(w, "Пользователь не найден в контексте", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Вы получили доступ! Ваш ID: %v", userID)))
}
