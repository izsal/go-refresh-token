package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/izsal/go-refresh-token/database"
	"github.com/izsal/go-refresh-token/models"
	"github.com/izsal/go-refresh-token/utils"

	"github.com/gorilla/mux"
)

func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	database.DB.Find(&items)
	response := utils.ItemsResponse{
		Status:  true,
		Message: "Items retrieved successfully",
		Items:   items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var item models.Item
	result := database.DB.First(&item, id)
	if result.Error != nil {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Item not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := utils.ItemResponse{
		Status:  true,
		Message: "Item retrieved successfully",
		Item:    item,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)

	// Validate input
	if item.Name == "" || item.Price <= 0 {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Name and price are required",
		}
		http.Error(w, "Invalid input", http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	database.DB.Create(&item)
	response := utils.ItemResponse{
		Status:  true,
		Message: "Item created successfully",
		Item:    item,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var item models.Item
	result := database.DB.First(&item, id)
	if result.Error != nil {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Item not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Body Json
	json.NewDecoder(r.Body).Decode(&item)

	// Validate input
	if item.Name == "" || item.Price <= 0 {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Name and price are required",
		}
		http.Error(w, "Invalid input", http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	database.DB.Save(&item)

	response := utils.ItemResponse{
		Status:  true,
		Message: "Item updated successfully",
		Item:    item,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 { // Validasi untuk memastikan id adalah angka positif
		response := utils.GenericResponse{
			Status:  false,
			Message: "Invalid ID",
		}
		http.Error(w, "Invalid input", http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var item models.Item
	result := database.DB.First(&item, id)
	if result.Error != nil {
		// Mengembalikan response JSON saat item tidak ditemukan
		response := utils.GenericResponse{
			Status:  false,
			Message: "Item Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hapus item
	result = database.DB.Delete(&item, id)
	if result.Error != nil {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Error deleting item",
		}
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
