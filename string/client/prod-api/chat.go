package prod_api

import (
	"string_um/string/models"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetChats(params map[string]interface{}) ([]models.Chat, error) {
	// Filter the Chats slice based on the parameters
	var filteredChats []models.Chat
	result := Database.Where(params).Find(&filteredChats)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredChats, nil
}

// Handler function to handle GET requests directly to database
func GetChat(id uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	result := Database.First(&chat, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &chat, nil
}

// Handler function to handle CREATE requests directly to database
func CreateChat(newChat models.Chat) (*models.Chat, error) {
	// Generate a new UUID for the chat if not provided
	if newChat.ID == uuid.Nil {
		newChat.ID = uuid.New()
	}

	// Add the new chat to the database
	result := Database.Create(&newChat)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newChat, nil
}

// Handler function to handle UPDATE requests directly to database
func UpdateChat(id uuid.UUID, partialChat map[string]interface{}) (*models.Chat, error) {
	// Update the chat with the provided ID
	var updatedChat models.Chat
	result := Database.Model(&updatedChat).Where("id = ?", id).Updates(partialChat)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedChat, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteChat(id uuid.UUID) error {
	// Remove the chat with the provided ID
	result := Database.Delete(&models.Chat{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
