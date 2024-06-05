package debug_api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"string_um/string/models"

	"errors"

	"github.com/google/uuid"
)

// Handler function to handle GET requests to /contactAddress endpoint
func GetContactAddresses(w http.ResponseWriter, r *http.Request) {
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

	// Build the filter conditions based on the query parameters
	filter := make(map[string]interface{})
	for key, values := range params {
		if len(values) > 0 {
			filter[key] = values[0]
		}
	}

	// Filter the ContactAddress entity based on the query parameters
	var filteredContactAddresses []models.ContactAddress
	result := Database.Where(filter).Find(&filteredContactAddresses)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredContactAddresses)
}

// Handler function to handle GET requests to /contactAddress/{id} endpoint
func GetContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Check if the ID is provided in the URL
	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	// Parse the ID from the URL
	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleError(w, err)
		return
	}

	var contactAddress models.ContactAddress
	result := Database.First(&contactAddress, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactAddress)
}

// Handler function to handle POST requests to /contactAddress/create endpoint
func CreateContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Parse JSON request body into a ContactAddress struct
	var newContactAddress models.ContactAddress
	if err := json.NewDecoder(r.Body).Decode(&newContactAddress); err != nil {
		HandleError(w, err)
		return
	}

	// Generate a new UUID for the contactAddres
	newContactAddress.ID = uuid.New()

	// Add the new contactAddres to the database
	result := Database.Create(&newContactAddress)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newContactAddress)
}

// Handler function to handle PUT requests to /contactAddress/update/{id} endpoint
func UpdateContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Check if the ID is provided in the URL
	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	// Parse the ID from the URL
	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleError(w, err)
		return
	}

	// Parse JSON request body into a map
	var partialContactAddress map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialContactAddress); err != nil {
		HandleError(w, err)
		return
	}

	// Update the contactAddress with the provided ID
	result := Database.Model(models.ContactAddress{}).Where("id = ?", uuid).Updates(partialContactAddress)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Retrieve the updated contactAddress
	var updatedContactAddress models.ContactAddress
	result = Database.First(&updatedContactAddress, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 200 and return the updated contactAddr
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedContactAddress)
}

// Handler function to handle DELETE requests to /ContactAddress endpoint
func DeleteContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	if r.PathValue("id") == "" {
		HandleError(w, errors.New("id is required"))
		return
	}

	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleError(w, err)
		return
	}

	// Remove the contactAddres with the provided ID
	result := Database.Delete(&models.ContactAddress{}, uuid)
	if result.RowsAffected == 0 {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
