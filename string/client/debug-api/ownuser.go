package debug_api

import (
	"encoding/json"
	"errors"
	"net/http"

	"string_um/string/models"
)

// Handler function to handle GET requests to /ownUser/{id} endpoint
func GetOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	var ownUser models.OwnUser
	result := Database.First(&ownUser)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ownUser)
}

// Handler function to handle POST requests to /ownUser/create endpoint
func CreateOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// If user already exists, return an error
	var ownUser models.OwnUser
	result := Database.Find(&ownUser)
	if result.RowsAffected > 0 {
		HandleError(w, errors.New("ownUser already exists"))
		return
	}

	// Parse JSON request body into a OwnUser struct
	var newOwnUser models.OwnUser
	if err := json.NewDecoder(r.Body).Decode(&newOwnUser); err != nil {
		HandleError(w, err)
		return
	}

	// Add the new ownUser to the database
	result = Database.Create(&newOwnUser)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOwnUser)
}

// Handler function to handle PUT requests to /ownUser/update/{id} endpoint
/* func UpdateOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the ID is provided in the URL
	if r.PathValue("id") == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Parse JSON request body into a map
	var partialOwnUser map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialOwnUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the ownUser with the provided ID
	var updatedOwnUser models.OwnUser
	result := Database.Model(&updatedOwnUser).Where("id = ?", r.PathValue("id")).Updates(partialOwnUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "OwnUser not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set status code to 200 and return the updated message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedOwnUser)
} */

// Handler function to handle DELETE requests to /ownUser endpoint
func DeleteOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Remove ownUser
	result := Database.Delete(&models.OwnUser{})
	if result.RowsAffected == 0 {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
