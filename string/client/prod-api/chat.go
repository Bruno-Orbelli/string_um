package prod_api

import (
	"string_um/string/entities"

	"github.com/google/uuid"
)

// Handler function to handle GET requests directly to database
func GetChats(params map[string]interface{}) ([]entities.Chat, error) {
	// Filter the Chats slice based on the parameters
	var filteredChats []entities.Chat
	result := Database.Preload("Messages").Where(params).Find(&filteredChats)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredChats, nil
}

// Handler function to handle GET requests directly to database
func GetChat(id uuid.UUID) (*entities.Chat, error) {
	var chat entities.Chat
	result := Database.Preload("Messages").First(&chat, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &chat, nil
}

// Handler function to handle CREATE requests directly to database
func CreateChat(newChat entities.Chat) (*entities.Chat, error) {
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
func UpdateChat(id uuid.UUID, partialChat map[string]interface{}) (*entities.Chat, error) {
	// Update the chat with the provided ID
	result := Database.Preload("Messages").Model(entities.Chat{}).Where("id = ?", id).Updates(partialChat)
	if result.Error != nil {
		return nil, result.Error
	}

	// Retrieve the updated chat
	var updatedChat entities.Chat
	result = Database.Preload("Messages").First(&updatedChat, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedChat, nil
}

// Handler function to handle DELETE requests directly to database
func DeleteChat(id uuid.UUID) error {
	// Remove the chat with the provided ID
	result := Database.Delete(&entities.Chat{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
