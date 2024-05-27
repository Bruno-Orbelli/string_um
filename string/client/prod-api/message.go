package prod_api

import (
	"string_um/string/models"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetMessages(params map[string]interface{}) ([]models.Message, error) {
	// Filter the Messages slice based on the parameters
	var filteredMessages []models.Message
	result := Database.Where(params).Find(&filteredMessages)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredMessages, nil
}

// Handler function to handle GET requests directly to database
func GetMessage(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	result := Database.First(&message, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &message, nil
}

// Handler function to handle CREATE requests directly to database
func CreateMessage(newMessage models.Message) (*models.Message, error) {
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
func UpdateMessage(id uuid.UUID, partialMessage map[string]interface{}) (*models.Message, error) {
	// Update the Message with the provided ID
	result := Database.Model(models.Message{}).Where("id = ?", id).Updates(partialMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated message
	var updatedMessage models.Message
	result = Database.First(&updatedMessage, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedMessage, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteMessage(id uuid.UUID) error {
	// Remove the Message with the provided ID
	result := Database.Delete(&models.Message{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
