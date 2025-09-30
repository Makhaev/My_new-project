package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"main.go/db"
	"main.go/models"
	"main.go/utils"
)

func Register(c *gin.Context) {
	var req struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// вот здесь сравнимаем логин
	var user models.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// здесь будем сравнивать пароль который у нас есть в базе данных и сырой пароль

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateJWT(user.ID)
	// здесь отправляем токен пользователю
	c.JSON(http.StatusOK, gin.H{"access": token})

}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}
