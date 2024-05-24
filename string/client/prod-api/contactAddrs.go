package prod_api

import (
	"string_um/string/models"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetContactAddresses(params map[string]interface{}) ([]models.ContactAddress, error) {
	// Filter the ContactAddresses slice based on the parameters
	var filteredContactAddresses []models.ContactAddress
	result := Database.Where(params).Find(&filteredContactAddresses)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredContactAddresses, nil
}

// Handler function to handle GET requests directly to database
func GetContactAddress(id uuid.UUID) (*models.ContactAddress, error) {
	var contactAddress models.ContactAddress
	result := Database.First(&contactAddress, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &contactAddress, nil
}

// Handler function to handle CREATE requests directly to database
func CreateContactAddress(newContactAddress models.ContactAddress) (*models.ContactAddress, error) {
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
func UpdateContactAddress(id uuid.UUID, partialContactAddress map[string]interface{}) (*models.ContactAddress, error) {
	// Update the contactAddress with the provided ID
	var updatedContactAddress models.ContactAddress
	result := Database.Model(&updatedContactAddress).Where("id = ?", id).Updates(partialContactAddress)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedContactAddress, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteContactAddress(id uuid.UUID) error {
	// Remove the contactAddress with the provided ID
	result := Database.Delete(&models.ContactAddress{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
