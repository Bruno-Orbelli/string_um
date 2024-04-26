package database_api

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	filteredMessages := Messages
	if r.URL.Query().Get("chatID") != "" {
		uuid, err := uuid.Parse(r.URL.Query().Get("chatID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, message := range filteredMessages {
			if message.Chat.ID != uuid {
				filteredMessages = append(filteredMessages[:i], filteredMessages[i+1:]...)
			}
		}
	}
	if r.URL.Query().Get("alreadySent") != "" {
		alreadySent, err := strconv.ParseBool(r.URL.Query().Get("alreadySent"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, message := range filteredMessages {
			if message.AlredySent != alreadySent {
				filteredMessages = append(filteredMessages[:i], filteredMessages[i+1:]...)
			}
		}
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

// Handler function to handle PUT requests to /messages endpoint
func UpdateMessageToSent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the ID is provided in the URL
	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the chat in the messages slice
	for _, message := range Messages {
		if message.ID == uuid {
			message.AlredySent = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Message not found", http.StatusNotFound)
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
