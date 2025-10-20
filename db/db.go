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

	// Таблица sms_codes
	createdSMSCodes := `
	CREATE TABLE IF NOT EXISTS sms_codes (
		id SERIAL PRIMARY KEY,
		phone TEXT NOT NULL,
		code TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP
	);`
	_, err = DB.Exec(createdSMSCodes)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы sms_codes: %v", err)
	}

	// Таблица refresh_tokens
	createdRefreshTokens := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		phone TEXT NOT NULL,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(createdRefreshTokens)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы refresh_tokens: %v", err)
	}

	// Таблица categories
	createdCategories := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		image TEXT,
		slug TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP
	);`
	_, err = DB.Exec(createdCategories)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы categories: %v", err)
	}

	// Таблица users
	createdUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
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
	_, err = DB.Exec(createdUsers)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}

	fmt.Println("База данных успешно подключена")
}
