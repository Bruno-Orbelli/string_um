package prod_api

import (
	"string_um/string/entities"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetContactAddresses(params map[string]interface{}) ([]entities.ContactAddress, error) {
	// Filter the ContactAddresses slice based on the parameters
	var filteredContactAddresses []entities.ContactAddress
	result := Database.Where(params).Find(&filteredContactAddresses)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredContactAddresses, nil
}

// Handler function to handle GET requests directly to database
func GetContactAddress(id uuid.UUID) (*entities.ContactAddress, error) {
	var contactAddress entities.ContactAddress
	result := Database.First(&contactAddress, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &contactAddress, nil
}

// Handler function to handle CREATE requests directly to database
func CreateContactAddress(newContactAddress entities.ContactAddress) (*entities.ContactAddress, error) {
	// Generate a new UUID for the contactAddress if not provided
	if newContactAddress.ID == uuid.Nil {
		newContactAddress.ID = uuid.New()
	}

	// Add the new contactAddress to the database
	result := Database.Create(&newContactAddress)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newContactAddress, nil
}

// Handler function to handle UPDATE requests directly to database
func UpdateContactAddress(id uuid.UUID, partialContactAddress map[string]interface{}) (*entities.ContactAddress, error) {
	// Update the contactAddress with the provided ID
	result := Database.Model(entities.ContactAddress{}).Where("id = ?", id).Updates(partialContactAddress)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated contactAddress
	var updatedContactAddress entities.ContactAddress
	result = Database.First(&updatedContactAddress, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedContactAddress, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteContactAddress(id uuid.UUID) error {
	// Remove the contactAddress with the provided ID
	result := Database.Delete(&entities.ContactAddress{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
