package models

import (
	"log"

	"main.go/db"
)

type UserStoreInfo struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	StoreName        string `json:"store_name"`
	StoreImage       string `json:"store_image"`
	StoreAddress     string `json:"store_address"`
	StorePhone       string `json:"store_phone"`
	StoreCode        string `json:"store_code"`
	ManagerName      string `json:"manager_name"`
	ManagerPhone     string `json:"manager_phone"`
	RemainingDebt    string `json:"remaining_debt"` // или float64, если нужна арифметика
	FavoriteProducts []int  `json:"favorite_products"`
}

func CreateUsers() {
	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phone TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    store_name TEXT,
    store_image TEXT,
    store_address TEXT,
    store_phone TEXT,
    store_code TEXT,
    manager_name TEXT,
    manager_phone TEXT,
    remaining_debt TEXT
);`

	_, err := db.DB.Exec(createUsersTable)

	if err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}

}

func (u *UserStoreInfo) CreateUser() error {
	query := `
	INSERT INTO users (
		phone, username, first_name, last_name, store_name, store_image,
		store_address, store_phone, store_code, manager_name, manager_phone, remaining_debt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.DB.Exec(query,
		u.StorePhone, // вместо u.Username
		u.Username,
		u.FirstName,
		u.LastName,
		u.StoreName,
		u.StoreImage,
		u.StoreAddress,
		u.StorePhone,
		u.StoreCode,
		u.ManagerName,
		u.ManagerPhone,
		u.RemainingDebt,
	)

	return err
}

func GetUserByPhone(phone string) (*UserStoreInfo, error) {
	query := `
	SELECT id, phone, username, first_name, last_name, store_name, store_image,
	       store_address, store_phone, store_code, manager_name, manager_phone, remaining_debt
	FROM users WHERE phone = ? LIMIT 1
	`

	row := db.DB.QueryRow(query, phone)

	var user UserStoreInfo
	err := row.Scan(
		&user.ID,
		&user.StorePhone, // сюда phone
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.StoreName,
		&user.StoreImage,
		&user.StoreAddress,
		&user.StorePhone,
		&user.StoreCode,
		&user.ManagerName,
		&user.ManagerPhone,
		&user.RemainingDebt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
