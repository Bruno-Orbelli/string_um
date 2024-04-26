package database_api

import (
	"encoding/json"
	"net/http"

	"string_um/string/models"
)

var OwnUsers []models.OwnUser

func GetOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(OwnUsers) == 0 {
		// Set status code to 404 (Not Found)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set content type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Marshal users slice to JSON
	json.NewEncoder(w).Encode(OwnUsers[0])
}

// Handler function to handle POST requests to /users endpoint
func CreateOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into a User struct
	var newOwnUser models.OwnUser
	err := json.NewDecoder(r.Body).Decode(&newOwnUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the new user to the users slice
	OwnUsers = append(OwnUsers, newOwnUser)

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
}

// Handler function to handle DELETE requests to /users endpoint
func DeleteOwnUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Remove the last user from the users slice
	OwnUsers = OwnUsers[:len(OwnUsers)-1]

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
