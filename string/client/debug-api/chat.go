package debug_api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"string_um/string/models"

	"errors"

	"github.com/google/uuid"
)

// Handler function to handle GET requests to /chats endpoint
func GetChats(w http.ResponseWriter, r *http.Request) {
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

	// Filter the Chats slice based on the query parameters
	var filteredChats []models.Chat
	result := Database.Preload("Messages").Where(filter).Find(&filteredChats)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredChats)
}

// Handler function to handle GET requests to /chats/{id} endpoint
func GetChat(w http.ResponseWriter, r *http.Request) {
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

	var chat models.Chat
	result := Database.Preload("Messages").First(&chat, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// Handler function to handle POST requests to /chats/create endpoint
func CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New("method not allowed"))
		return
	}

	// Parse JSON request body into a Chat struct
	var newChat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&newChat); err != nil {
		HandleError(w, err)
		return
	}

	// Generate a new UUID for the chat if not provided
	if newChat.ID == uuid.Nil {
		newChat.ID = uuid.New()
	}

	// Add the new chat to the database
	result := Database.Create(&newChat)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newChat)
}

// Handler function to handle PUT requests to /chats/update/{id} endpoint
func UpdateChat(w http.ResponseWriter, r *http.Request) {
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
	var partialChat map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialChat); err != nil {
		HandleError(w, err)
		return
	}

	// Update the chat with the provided ID
	result := Database.Preload("Messages").Model(models.Chat{}).Where("id = ?", uuid).Updates(partialChat)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Retrieve the updated chat
	var updatedChat models.Chat
	result = Database.Preload("Messages").First(&updatedChat, "id = ?", uuid)
	if result.Error != nil {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 200 and return the updated chat
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedChat)
}

// Handler function to handle DELETE requests to /Chats endpoint
func DeleteChat(w http.ResponseWriter, r *http.Request) {
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

	// Remove the chat with the provided ID
	result := Database.Delete(&models.Chat{}, uuid)
	if result.RowsAffected == 0 {
		HandleError(w, result.Error)
		return
	}

	// Set status code to 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
