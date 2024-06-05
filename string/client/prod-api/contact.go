package prod_api

import (
	"string_um/string/models"
)

// Handler function to handle GET requests directly to database
func GetContacts(params map[string]interface{}) ([]models.Contact, error) {
	// Filter the Contacts slice based on the parameters
	var filteredContacts []models.Contact
	result := Database.Where(params).Find(&filteredContacts)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredContacts, nil
}

// Handler function to handle GET requests directly to database
func GetContact(id string) (*models.Contact, error) {
	var contact models.Contact
	result := Database.Preload("ContactAddresses").Preload("Chat").First(&contact, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &contact, nil
}

// Handler function to handle CREATE requests directly to database
func CreateContact(newContact models.Contact) (*models.Contact, error) {
	// Add the new contact to the database
	result := Database.Create(&newContact)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newContact, nil
}

// Handler function to handle UPDATE requests directly to database
func UpdateContact(id string, partialContact map[string]interface{}) (*models.Contact, error) {
	// Update the contact with the provided ID
	result := Database.Preload("ContactAddresses").Preload("Chat").Model(models.Contact{}).Where("id = ?", id).Updates(partialContact)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated contact
	var updatedContact models.Contact
	result = Database.Preload("ContactAddresses").Preload("Chat").First(&updatedContact, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedContact, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteContact(id string) error {
	// Remove the contact with the provided ID
	result := Database.Delete(&models.Contact{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
