package models

import "gorm.io/gorm"

// Шаг 2. Создание модели
type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	PasswordHash string
	FirstName    string
	LastName     string
	StoreName    string
	StoreAddress string
	StorePhone   string
	StoreCode    string
	ManagerName  string
	ManagerPhone string
}
