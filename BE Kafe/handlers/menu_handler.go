package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"be_kafe/config"
	"be_kafe/models"
	"be_kafe/utils"

	"github.com/gorilla/mux"
)

// CreateMenuHandler handles menu creation
func CreateMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var menu models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if menu.Nama == "" || menu.Foto == "" || menu.Harga <= 0 || menu.Deskripsi == "" {
		http.Error(w, "Semua field harus diisi dan harga harus lebih dari 0", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(
		"INSERT INTO menus (nama, foto, harga, deskripsi) VALUES (?, ?, ?, ?)",
		menu.Nama, menu.Foto, menu.Harga, menu.Deskripsi,
	)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Gagal menambah menu", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	menu.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(menu); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// UpdateMenuHandler handles menu updates
func UpdateMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var menu models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if menu.Nama == "" || menu.Foto == "" || menu.Harga <= 0 || menu.Deskripsi == "" {
		http.Error(w, "Semua field harus diisi dan harga harus lebih dari 0", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(
		"UPDATE menus SET nama=?, foto=?, harga=?, deskripsi=? WHERE id=?",
		menu.Nama, menu.Foto, menu.Harga, menu.Deskripsi, id,
	)
	if err != nil {
		log.Printf("Database update error: %v", err)
		http.Error(w, "Gagal mengedit menu", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Menu tidak ditemukan", http.StatusNotFound)
		return
	}

	menu.ID = utils.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(menu); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// DeleteMenuHandler handles menu deletion
func DeleteMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DB.Exec("DELETE FROM menus WHERE id=?", id)
	if err != nil {
		log.Printf("Database delete error: %v", err)
		http.Error(w, "Gagal menghapus menu", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Menu tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Menu berhasil dihapus"})
}
