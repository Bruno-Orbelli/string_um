package prod_api

import (
	"errors"

	"string_um/string/entities"
)

// Handler function to handle GET requests directly to database
func GetOwnUser() (*entities.OwnUser, error) {
	var ownUser entities.OwnUser
	result := Database.First(&ownUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ownUser, nil
}

// Handler function to handle CREATE requests directly to database
func CreateOwnUser(newOwnUser entities.OwnUser) (*entities.OwnUser, error) {
	// If user already exists, return an error
	var ownUser entities.OwnUser
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
func UpdateOwnUser(id string, partialUser map[string]interface{}) (*entities.OwnUser, error) {
	// Update the OwnUser with the provided partial data
	result := Database.Model(entities.OwnUser{}).Where("id = ?", id).Updates(partialUser)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated ownUser
	var updatedOwnUser entities.OwnUser
	result = Database.First(&updatedOwnUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedOwnUser, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteOwnUser() error {
	// Remove ownUser
	result := Database.Delete(&entities.OwnUser{})
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
