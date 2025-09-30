package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}

	return 0, err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		userID, err := ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
