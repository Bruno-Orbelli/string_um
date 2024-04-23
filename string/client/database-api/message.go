package client

import (
	"encoding/json"
	"net/http"
	"time"

	"string_um/string/models"

	"github.com/google/uuid"
)

var Messages []models.Message

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var filteredMessages []models.Message
	if r.URL.Query().Get("chat_id") != "" {
		uuid, err := uuid.Parse(r.URL.Query().Get("chat_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, message := range Messages {
			if message.Chat.ID == uuid {
				filteredMessages = append(filteredMessages, message)
			}
		}
	} else {
		copy(filteredMessages, Messages)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredMessages)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, message := range Messages {
		if message.ID == uuid {
			// Set content type header to application/json
			w.Header().Set("Content-Type", "application/json")

			// Marshal messages slice to JSON
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}

// Handler function to handle POST requests to /messages endpoint
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into a Message struct
	var newMessage models.Message
	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newMessage.SentAt = time.Now()
	if newMessage.SentBy.ID == OwnUsers[0].ID {
		newMessage.AlredySent = false
	}

	// Add the new chat to the messages slice
	Messages = append(Messages, newMessage)

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
}

// Handler function to handle DELETE requests to /messages endpoint
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	// Remove the last chat from the messages slice
	for i, chat := range Messages {
		if chat.ID == uuid {
			Messages = append(Messages[:i], Messages[i+1:]...)

			// Set status code to 204 (No Content)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Set status code to 404 (Not Found)
	w.WriteHeader(http.StatusNotFound)
}
