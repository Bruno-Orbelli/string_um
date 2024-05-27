package funcs

import (
	"errors"
	"fmt"
	"os"
	prod_api "string_um/string/client/prod-api"
	"string_um/string/main/encryption"
	"string_um/string/models"
	"time"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"gorm.io/gorm"
)

func CreateOwnUser() (*models.OwnUser, *models.Contact, error) {
	// fmt.Println("Creating own user...")

	// Generate private key
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Marshal private key to bytes
	marshalledKey, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	// Get peer ID from private key
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	// Create own user
	ownUser := models.OwnUser{
		ID:         id.String(),
		PrivateKey: marshalledKey,
	}

	// Persist own user
	createdOwnUser, err := prod_api.CreateOwnUser(ownUser)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create own user: %w", err)
	}

	// Create own user contact
	ownUserContact := models.Contact{
		ID:   createdOwnUser.ID,
		Name: "Me",
	}

	// Persist own user contact
	createdOwnContact, err := prod_api.CreateContact(ownUserContact)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create own user contact: %w", err)
	}

	return createdOwnUser, createdOwnContact, nil
}

func GetOwnUserIfExists() (*models.OwnUser, *models.Contact, error) {
	// Check if own user exists
	ownUser, err := prod_api.GetOwnUser()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	if ownUser != nil { // Own user exists
		// Retrieve own user contact
		ownUserContact, err := prod_api.GetContact(ownUser.ID)
		if err != nil {
			return nil, nil, err
		}
		return ownUser, ownUserContact, nil
	} else { // Own user does not exist
		return CreateOwnUser()
	}
}

func AddMessageToBeSent(chatID uuid.UUID, senderID string, body string) error {
	message := models.Message{
		ID:          uuid.New(),
		ChatID:      chatID,
		AlreadySent: false,
		SentByID:    senderID,
		SentAt:      time.Now(),
		Message:     body,
	}
	createdMessage, err := prod_api.CreateMessage(message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	fmt.Printf("Created message with ID=%s.\n", createdMessage.ID)
	return nil
}

func Register(password string) error {
	prod_api.RunDatabaseAPI()
	if _, _, err := CreateOwnUser(); err != nil {
		return err
	} else {
		encryption.EncryptFile("test.db", "en_test.db", password, "en_salt.txt")
		os.Remove("test.db")
	}
	return nil
}

func Login(password string) error {
	if err := encryption.DecryptFile("en_test.db", "test.db", password, "en_salt.txt"); err != nil {
		return err
	}
	return nil
}

func GetChatsAndInfoWithMessages() ([]models.ChatDTO, error) {
	chats, err := prod_api.GetChats(nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	var chatDTOsWithMessages []models.ChatDTO
	for _, chat := range chats {
		if len(chat.Messages) > 0 {
			contact, err := prod_api.GetContact(chat.ContactID)
			if err != nil {
				return nil, fmt.Errorf("failed to get contact: %w", err)
			}
			// Create chat DTO with messages
			chatDTO := models.ChatDTO{
				ID:          chat.ID,
				ContactName: contact.Name,
				Messages:    chat.Messages,
			}
			chatDTOsWithMessages = append(chatDTOsWithMessages, chatDTO)
		}
	}

	return chatDTOsWithMessages, nil
}

/*func InitDatabase(password string) error {
	// Check if database file exists and
	if _, err := os.Stat("en_test.db"); err != nil {
		if os.IsNotExist(err) {
			// Create the database file.
			_, err := os.Create("en_test.db")
			if err != nil {
				return err
			}
			if err := encryption.EncryptFile("en_test.db", "en_test.db", password, "en_salt.txt"); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if err := encryption.DecryptFile("en_test.db", "test.db", password, "en_salt.txt"); err != nil {
		return err
	}

	return nil
}*/

func GetChatWithContact(contactID string) (*models.Chat, error) {
	var chats []models.Chat
	var destContact *models.Contact

	// Get contact to send message to.
	destContact, err := prod_api.GetContact(contactID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	// Get chat with contact.
	params := map[string]interface{}{"contact_id": destContact.ID}
	chats, err = prod_api.GetChats(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	var chat models.Chat
	if len(chats) > 1 { // Unexpected number of chats, this should not happen.
		return nil, fmt.Errorf("fatal: unexpected number of chats for single contact: %s", destContact.Name)
	} else if len(chats) == 0 { // Chat does not exist, create it.
		fmt.Printf("Creating new chat with %s.\n", destContact.Name)
		chat = models.Chat{
			ContactID: destContact.ID,
		}
		createdChat, err := prod_api.CreateChat(chat)
		if err != nil {
			return nil, fmt.Errorf("failed to create chat: %w", err)
		}
		chat = *createdChat
	} else { // Chat exists, use it.
		chat = chats[0]
	}

	return &chat, nil
}
