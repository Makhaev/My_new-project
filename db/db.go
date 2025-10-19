package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Ошибка пинга базы данных: %v", err)
	}

	createdDB := `
	CREATE TABLE IF NOT EXISTS sms_codes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phone TEXT NOT NULL,
		code TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME
	);`

	_, err = DB.Exec(createdDB)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	fmt.Println("База данных успешно подключена")
	createdRefreshTable := `
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	phone TEXT NOT NULL,
	token TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	_, err = DB.Exec(createdRefreshTable)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы refresh_tokens: %v", err)
	}

	// Таблица для карточек

	createdCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT,
    image TEXT,
    slug TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
	)
	`

	_, err = DB.Exec(createdCategoriesTable)

	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

}
