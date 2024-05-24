package prod_api

import (
	"errors"

	"string_um/string/models"
)

// Handler function to handle GET requests directly to database
func GetOwnUser() (*models.OwnUser, error) {
	var ownUser models.OwnUser
	result := Database.First(&ownUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ownUser, nil
}

// Handler function to handle CREATE requests directly to database
func CreateOwnUser(newOwnUser models.OwnUser) (*models.OwnUser, error) {
	// If user already exists, return an error
	var ownUser models.OwnUser
	result := Database.Find(&ownUser)
	if result.RowsAffected > 0 {
		return nil, errors.New("ownUser already exists")
	}

	// Add the new ownUser to the database
	result = Database.Create(&newOwnUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newOwnUser, nil
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

// Handler function to handle DELETE requests directly to database
func DeleteOwnUser() error {
	// Remove ownUser
	result := Database.Delete(&models.OwnUser{})
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
