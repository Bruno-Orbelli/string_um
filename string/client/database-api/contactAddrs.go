package client

import (
	"encoding/json"
	"net/http"
	"time"

	"string_um/string/models"

	"github.com/google/uuid"
)

var contactAddrs []models.ContactAddress

func GetContactAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactAddrs)
}

func GetContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, contactAddr := range contactAddrs {
		if contactAddr.ID == uuid {
			// Set content type header to application/json
			w.Header().Set("Content-Type", "application/json")

			// Marshal users slice to JSON
			json.NewEncoder(w).Encode(contactAddr)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}

// Handler function to handle POST requests to /contact-addresses endpoint
func CreateContactAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into a ContactAddress struct
	var newContactAddr models.ContactAddress
	err := json.NewDecoder(r.Body).Decode(&newContactAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newContactAddr.ObservedAt = time.Now()

	// Add the new addr to the addrs slice
	contactAddrs = append(contactAddrs, newContactAddr)

	// Set status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
}

// Handler function to handle DELETE requests to /contact-addresses endpoint
func DeleteContactAddress(w http.ResponseWriter, r *http.Request) {
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

	// Remove the last addr from the addrs slice
	for i, contactAddr := range contactAddrs {
		if contactAddr.ID == uuid {
			contactAddrs = append(contactAddrs[:i], contactAddrs[i+1:]...)

			// Set status code to 204 (No Content)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Set status code to 404 (Not Found)
	w.WriteHeader(http.StatusNotFound)
}
