package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"main.go/models"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	// можно валидировать имя
	if c.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	// генерируем slug, если пустой (например простая замена пробелов)
	if c.Slug == "" {
		c.Slug = generateSlug(c.Name)
	}
	id, err := c.Create()
	if err != nil {
		http.Error(w, "Ошибка создания категории", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	list, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Ошибка получения категорий", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	c, err := models.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Ошибка получения категории", http.StatusInternalServerError)
		return
	}
	if c == nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	c.ID = id
	if c.Slug == "" {
		c.Slug = generateSlug(c.Name)
	}
	if err := models.UpdateCategory(&c); err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCategory(id); err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// можно удалить лишние символы и т.д.
	return slug
}
