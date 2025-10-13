package models

import (
	"database/sql"
	"time"

	"main.go/db"
)

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (c *Category) Create() (int64, error) {
	q := `INSERT INTO categories (name, description, image, slug, created_at) VALUES (?, ?, ?, ?, ?)`
	res, err := db.DB.Exec(q, c.Name, c.Description, c.Image, c.Slug, time.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetCategoryByID(id int) (*Category, error) {
	q := `SELECT id, name, description, image, slug, created_at, updated_at FROM categories WHERE id = ? LIMIT 1`
	row := db.DB.QueryRow(q, id)
	var c Category
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Image, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func GetAllCategories() ([]Category, error) {
	q := `SELECT id, name, description, image, slug, created_at, updated_at FROM categories ORDER BY name`
	rows, err := db.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Image, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func UpdateCategory(c *Category) error {
	q := `UPDATE categories SET name = ?, description = ?, image = ?, slug = ?, updated_at = ? WHERE id = ?`
	_, err := db.DB.Exec(q, c.Name, c.Description, c.Image, c.Slug, time.Now(), c.ID)
	return err
}

func DeleteCategory(id int) error {
	_, err := db.DB.Exec(`DELETE FROM categories WHERE id = ?`, id)
	return err
}
