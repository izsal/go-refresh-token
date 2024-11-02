package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/izsal/go-refresh-token/database"
	"github.com/izsal/go-refresh-token/logger"
	"github.com/izsal/go-refresh-token/models"
	"github.com/izsal/go-refresh-token/utils"

	"log"

	"github.com/gorilla/mux"
)

var logFiles *os.File

func init() {
	logFiles = logger.InitLogger()
	log.SetOutput(logFiles) // Set logger to write to logFile
}

// Log a user action related to item operations
func logItemAction(action string, id int, status string) {
	log.Printf("Action: %s | Item ID: %d | Status: %s\n", action, id, status)
}

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

	logItemAction("GetItems", 0, "Success")
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

		logItemAction("GetItem", id, "Failed - Item not found")
		return
	}

	response := utils.ItemResponse{
		Status:  true,
		Message: "Item retrieved successfully",
		Item:    item,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logItemAction("GetItem", id, "Success")
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

		logItemAction("CreateItem", 0, "Failed - Invalid input")
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

	logItemAction("CreateItem", int(response.Item.ID), "Success")
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

		logItemAction("UpdateItem", id, "Failed - Item not found")
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

		logItemAction("UpdateItem", id, "Failed - Invalid input")
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

	logItemAction("UpdateItem", id, "Success")
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 { // Validate to ensure id is a positive number
		response := utils.GenericResponse{
			Status:  false,
			Message: "Invalid ID",
		}
		http.Error(w, "Invalid input", http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)

		logItemAction("DeleteItem", 0, "Failed - Invalid ID")
		return
	}

	var item models.Item
	result := database.DB.First(&item, id)
	if result.Error != nil {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Item Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)

		logItemAction("DeleteItem", id, "Failed - Item not found")
		return
	}

	// Delete item
	result = database.DB.Delete(&item, id)
	if result.Error != nil {
		response := utils.GenericResponse{
			Status:  false,
			Message: "Error deleting item",
		}
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		logItemAction("DeleteItem", id, "Failed - Deletion error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logItemAction("DeleteItem", id, "Success")
}
