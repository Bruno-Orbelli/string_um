package database_api

import (
	"encoding/json"
	"net/http"

	"string_um/string/models"

	"github.com/google/uuid"
)

var Chats []models.Chat

func GetChats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filteredChats := Chats
	if r.URL.Query().Get("contactID") != "" {
		id := r.URL.Query().Get("contactID")
		for i, chat := range filteredChats {
			if chat.Contact.ID != id {
				filteredChats = append(filteredChats[:i], filteredChats[i+1:]...)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredChats)
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, chat := range Chats {
		if chat.ID == uuid {
			// Set content type header to application/json
			w.Header().Set("Content-Type", "application/json")

			// Marshal Chats slice to JSON
			json.NewEncoder(w).Encode(chat)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}

// Handler function to handle POST requests to /Chats endpoint
func CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into a Chat struct
	var newChat models.Chat
	err := json.NewDecoder(r.Body).Decode(&newChat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the new chat to the Chats slice
	Chats = append(Chats, newChat)

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
}

// Handler function to handle DELETE requests to /Chats endpoint
func DeleteChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.PathValue("id") == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Remove the last chat from the Chats slice
	for i, chat := range Chats {
		if chat.ID == uuid {
			Chats = append(Chats[:i], Chats[i+1:]...)

			// Set status code to 204 (No Content)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Set status code to 404 (Not Found)
	w.WriteHeader(http.StatusNotFound)
}
