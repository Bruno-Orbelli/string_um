package funcs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	prod_api "string_um/string/client/prod-api"
	"string_um/string/main/encryption"
	"string_um/string/models"
	"string_um/string/networking/node"
	"time"

	"github.com/google/uuid"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"gorm.io/gorm"
)

func CreateOwnUser(passwordHash string) (*models.OwnUser, *models.Contact, error) {
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
		ID:           id.String(),
		PasswordHash: passwordHash,
		PrivateKey:   marshalledKey,
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

func AddContactAddressesForUnknownContacts(host host.Host) error {
	contacts, err := prod_api.GetContacts(nil)
	if err != nil {
		return fmt.Errorf("failed to get contacts: %w", err)
	}

	for _, contact := range contacts {
		if contact.ID == contact.Name {
			err = node.AddKnownAddressesForContact(host, contact.ID)
			if err != nil {
				return fmt.Errorf("failed to add known addresses for contact: %w", err)
			}
		}
	}

	return nil
}

func StartHost(ctx context.Context, priv crypto.PrivKey) (host.Host, *dht.IpfsDHT, error) {
	// Get key if it does not exist
	var privKey crypto.PrivKey
	var err error
	if priv == nil {
		privKey, _, err = crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			return nil, nil, err
		}
	} else {
		privKey = priv
	}

	// Try the default multiaddresses if the ports are available
	var listenAddrs []multiaddr.Multiaddr
	portNum := 40000
	for {
		// Check port availability
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
		if err != nil {
			portNum++
			continue
		}
		ln.Close()

		listenAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip6/::/tcp/%d", portNum))
		if err != nil {
			return nil, nil, err
		}
		listenAddrs = append(listenAddrs, listenAddr)
		break
	}

	host, dht, err := node.CreateNewNode(
		ctx,
		privKey,
		listenAddrs,
		make([]multiaddr.Multiaddr, 0),
		make([]multiaddr.Multiaddr, 0),
		"/chat/0.0.1",
	)
	if err != nil {
		return nil, nil, err
	}
	return host, dht, nil
}

func GetOwnUser() (*models.OwnUser, *models.Contact, error) {
	prod_api.RunDatabaseAPI()
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
		return nil, nil, nil
	}
}

func UpdateOwnUserHash(passwordHash string) error {
	prod_api.RunDatabaseAPI()
	ownUser, _, err := GetOwnUser()
	if err != nil {
		return err
	}

	params := map[string]interface{}{"password_hash": passwordHash}

	if _, err := prod_api.UpdateOwnUser(ownUser.ID, params); err != nil {
		return fmt.Errorf("failed to update own user: %w", err)
	}

	return nil
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
	_, err := prod_api.CreateMessage(message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// fmt.Printf("Created message with ID=%s.\n", createdMessage.ID)
	return nil
}

func Register(password string) error {
	prod_api.RunDatabaseAPI()
	salt, err := encryption.GenerateSalt("en_salt.txt")
	if err != nil {
		return err
	}
	hash, err := encryption.HashPassword(password, []byte(salt))
	if err != nil {
		return err
	}
	if _, _, err := CreateOwnUser(hash); err != nil {
		return err
	} else {
		encryption.EncryptFile("test.db", "en_test.db", hash, salt)
		os.Remove("test.db")
	}
	return nil
}

func Login(password string) error {
	// Get old salt
	oldSalt, err := encryption.GetSalt("en_salt.txt")
	if err != nil {
		return err
	}

	// Decrypt the database using the derived key
	if err := encryption.DecryptFile("en_test.db", "test.db", password, oldSalt); err != nil {
		return err
	}

	// Generate a new salt
	newSalt, err := encryption.GenerateSalt("en_salt.txt")
	if err != nil {
		return err
	}

	// Hash the password with the new salt
	newHash, err := encryption.HashPassword(password, []byte(newSalt))
	if err != nil {
		return err
	}

	// Update the user's hash and salt
	if err := UpdateOwnUserHash(newHash); err != nil {
		return err
	}

	return nil
}

func CloseDatabase() error {
	ownUser, _, err := GetOwnUser()
	if err != nil {
		return err
	}
	salt, err := encryption.GetSalt("en_salt.txt")
	if err != nil {
		return err
	}
	if err := encryption.EncryptFile("test.db", "en_test.db", ownUser.PasswordHash, salt); err != nil {
		return err
	}
	fileInfo, err := os.Stat("test.db")
	if err != nil {
		return err
	} else if fileInfo != nil {
		if err := os.Remove("test.db"); err != nil {
			return err
		}
	}
	return nil
}

func GetContactByName(name string) (*models.Contact, error) {
	contacts, err := prod_api.GetContacts(map[string]interface{}{"name": name})
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	if len(contacts) == 0 {
		return nil, nil
	}

	return &contacts[0], nil
}

func GetChatsAndInfo() ([]models.ChatDTO, error) {
	chats, err := prod_api.GetChats(nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	var chatDTOsWithMessages []models.ChatDTO
	for _, chat := range chats {
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
