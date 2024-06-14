package prod_api

import (
	"string_um/string/entities"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetMessages(params map[string]interface{}) ([]entities.Message, error) {
	// Filter the Messages slice based on the parameters
	var filteredMessages []entities.Message
	result := Database.Where(params).Find(&filteredMessages)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredMessages, nil
}

// Handler function to handle GET requests directly to database
func GetMessage(id uuid.UUID) (*entities.Message, error) {
	var message entities.Message
	result := Database.First(&message, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &message, nil
}

// Handler function to handle CREATE requests directly to database
func CreateMessage(newMessage entities.Message) (*entities.Message, error) {
	// Generate a new UUID for the Message if not provided
	if newMessage.ID == uuid.Nil {
		newMessage.ID = uuid.New()
	}

	// Add the new Message to the database
	result := Database.Create(&newMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newMessage, nil
}

// Handler function to handle UPDATE requests directly to database
func UpdateMessage(id uuid.UUID, partialMessage map[string]interface{}) (*entities.Message, error) {
	// Update the Message with the provided ID
	result := Database.Model(entities.Message{}).Where("id = ?", id).Updates(partialMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated message
	var updatedMessage entities.Message
	result = Database.First(&updatedMessage, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedMessage, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteMessage(id uuid.UUID) error {
	// Remove the Message with the provided ID
	result := Database.Delete(&entities.Message{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
