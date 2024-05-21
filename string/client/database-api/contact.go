package database_api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"string_um/string/models"
)

// Handler function to handle GET requests to /contacts endpoint
func GetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Get filtering parameters from the URL query
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Marshal the query parameters into a Contact struct
	jsonContact, err := json.Marshal(params)
	if err != nil {
		HandleError(w, err)
		return
	}
	var contact models.Contact
	if err := json.Unmarshal(jsonContact, &contact); err != nil {
		HandleError(w, err)
		return
	}

	// Filter contacts based on the query parameters
	var filteredContacts []models.Contact
	result := Database.Where(contact).Find(&filteredContacts)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredContacts)
}

// Handler function to handle GET requests to /contacts/{id} endpoint
func GetContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Check if the ID is provided in the URL
	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	var contact models.Contact
	result := Database.First(&contact, "id = ?", r.PathValue("id"))
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}

// Handler function to handle POST requests to /contacts/create endpoint
func CreateContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Parse JSON request body into a Contact struct
	var newContact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		HandleError(w, err)
		return
	}

	// Add the new contact to the database
	result := Database.Create(&newContact)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newContact)
}

// Handler function to handle PUT requests to /contacts/update/{id} endpoint
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Check if the ID is provided in the URL
	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	// Parse JSON request body into a map
	var partialContact map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialContact); err != nil {
		HandleError(w, err)
		return
	}

	// Update the contact with the provided ID
	var updatedContact models.Contact
	result := Database.Model(&updatedContact).Where("id = ?", r.PathValue("id")).Updates(partialContact)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 200 and return the updated contact
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedContact)
}

// Handler function to handle DELETE requests to /contacts endpoint
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	// Remove the contact with the provided ID
	result := Database.Delete(&models.Contact{}, r.PathValue("id"))
	if result.RowsAffected == 0 {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
