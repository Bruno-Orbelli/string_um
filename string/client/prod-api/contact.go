package prod_api

import (
	"string_um/string/entities"
)

// Handler function to handle GET requests directly to database
func GetContacts(params map[string]interface{}) ([]entities.Contact, error) {
	// Filter the Contacts slice based on the parameters
	var filteredContacts []entities.Contact
	result := Database.Preload("ContactAddresses").Preload("Chat").Where(params).Find(&filteredContacts)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredContacts, nil
}

// Handler function to handle GET requests directly to database
func GetContact(id string) (*entities.Contact, error) {
	var contact entities.Contact
	result := Database.Preload("ContactAddresses").Preload("Chat").First(&contact, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &contact, nil
}

// Handler function to handle CREATE requests directly to database
func CreateContact(newContact entities.Contact) (*entities.Contact, error) {
	// Add the new contact to the database
	result := Database.Create(&newContact)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newContact, nil
}

// Handler function to handle UPDATE requests directly to database
func UpdateContact(id string, partialContact map[string]interface{}) (*entities.Contact, error) {
	// Update the contact with the provided ID
	result := Database.Preload("ContactAddresses").Preload("Chat").Model(entities.Contact{}).Where("id = ?", id).Updates(partialContact)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated contact
	var updatedContact entities.Contact
	result = Database.Preload("ContactAddresses").Preload("Chat").First(&updatedContact, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedContact, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteContact(id string) error {
	// Remove the contact with the provided ID
	result := Database.Delete(&entities.Contact{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
