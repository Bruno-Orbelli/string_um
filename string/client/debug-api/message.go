package debug_api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"string_um/string/models"

	"errors"

	"github.com/google/uuid"
)

// Handler function to handle GET requests to /messages endpoint
func GetMessages(w http.ResponseWriter, r *http.Request) {
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

	// Filter the Messages based on the query parameters
	var filteredMessages []models.Message
	result := Database.Where(filter).Find(&filteredMessages)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredMessages)
}

// Handler function to handle GET requests to /messages/{id} endpoint
func GetMessage(w http.ResponseWriter, r *http.Request) {
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

	var message models.Message
	result := Database.First(&message, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// Handler function to handle POST requests to /messages/create endpoint
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Parse JSON request body into a Message struct
	var newMessage models.Message
	if err := json.NewDecoder(r.Body).Decode(&newMessage); err != nil {
		HandleError(w, err)
		return
	}

	// Generate a new UUID for the message
	newMessage.ID = uuid.New()

	// Add the new message to the database
	result := Database.Create(&newMessage)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newMessage)
}

// Handler function to handle PUT requests to /messages/update/{id} endpoint
func UpdateMessage(w http.ResponseWriter, r *http.Request) {
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
	var partialMessage map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialMessage); err != nil {
		HandleError(w, err)
		return
	}

	// Update the message with the provided ID
	result := Database.Model(models.Message{}).Where("id = ?", uuid).Updates(partialMessage)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Retrieve the updated message
	var updatedMessage models.Message
	result = Database.First(&updatedMessage, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 200 and return the updated message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMessage)
}

// Handler function to handle DELETE requests to /messages endpoint
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	// Remove the message with the provided ID
	result := Database.Delete(&models.Message{}, uuid)
	if result.RowsAffected == 0 {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
