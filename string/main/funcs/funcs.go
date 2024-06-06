package funcs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	prod_api "string_um/string/client/prod-api"
	"string_um/string/entities"
	"string_um/string/main/encryption"
	"string_um/string/main/images"
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

func CreateOwnUser(passwordHash string, imgEncodingHash string) (*entities.OwnUser, *entities.Contact, error) {
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
	ownUser := entities.OwnUser{
		ID:           id.String(),
		PasswordHash: passwordHash,
		EncodingHash: imgEncodingHash,
		PrivateKey:   marshalledKey,
	}

	// Persist own user
	createdOwnUser, err := prod_api.CreateOwnUser(ownUser)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create own user: %w", err)
	}

	// Create own user contact
	ownUserContact := entities.Contact{
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

func GetOwnUser() (*entities.OwnUser, *entities.Contact, error) {
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

func UpdateOwnUserHashes(passwordHash string, encodingHash string) error {
	prod_api.RunDatabaseAPI()
	ownUser, _, err := GetOwnUser()
	if err != nil {
		return err
	}

	params := map[string]interface{}{"password_hash": passwordHash, "encoding_hash": encodingHash}

	if _, err := prod_api.UpdateOwnUser(ownUser.ID, params); err != nil {
		return fmt.Errorf("failed to update own user: %w", err)
	}

	return nil
}

func AddMessageToBeSent(chatID uuid.UUID, senderID string, body string) error {
	message := entities.Message{
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

func CaptureMultipleImagesEncoding(numImg int) ([]float64, error) {
	// Capture images
	encodings := make([][]float32, numImg)
	for i := 0; i < numImg; i++ {
		img, err := images.CaptureImage()
		if err != nil {
			return nil, err
		}

		// Extract features
		encoding, err := images.ExtractFeatures(img)
		if err != nil {
			return nil, err
		}

		encodings[i] = encoding
	}

	// Average encodings
	avgEncoding := images.AverageEncodings(encodings)

	// Normalize average encoding
	normalizedEncoding := images.NormalizeEncoding(avgEncoding)

	return normalizedEncoding, nil
}

func CaptureSingleImageEncoding() ([]float64, error) {
	// Capture image
	img, err := images.CaptureImage()
	if err != nil {
		return nil, err
	}

	// Extract features
	encoding, err := images.ExtractFeatures(img)
	if err != nil {
		return nil, err
	}

	// Normalize encoding
	normalizedEncoding := images.NormalizeEncoding(encoding)

	return normalizedEncoding, nil
}

func Register(password string, normalizedEncoding []float64) error {
	prod_api.RunDatabaseAPI()
	salt, err := encryption.GenerateSalt("en_salt.txt")
	if err != nil {
		return err
	}
	passwordHash, err := encryption.HashInput(password, []byte(salt))
	if err != nil {
		return err
	}
	encodingHash, err := encryption.HashInput(normalizedEncoding, []byte(salt))
	if err != nil {
		return err
	}
	if _, _, err := CreateOwnUser(passwordHash, encodingHash); err != nil {
		return err
	} else {
		if err := encryption.EncryptFile("test.db", "en_test.db", encodingHash, salt); err != nil {
			return err
		}
		if err := encryption.EncryptFile("en_test.db", "en_test.db", passwordHash, salt); err != nil {
			return err
		}
		os.Remove("test.db")
	}
	return nil
}

func LoginFirstFactor(password string) error {
	// Get old salt
	oldSalt, err := encryption.GetSalt("en_salt.txt")
	if err != nil {
		return err
	}

	// Make the first decription using the password and the old salt
	if err := encryption.DecryptFile("en_test.db", "test.db", password, oldSalt); err != nil {
		return err
	}

	return nil
}

func LoginSecondFactor(password string, normalizedEncoding []float64) error {
	// Get old salt
	oldSalt, err := encryption.GetSalt("en_salt.txt")
	if err != nil {
		return err
	}

	// Make the second decryption using the encoding and the old salt
	if err := encryption.DecryptFile("test.db", "test.db", normalizedEncoding, oldSalt); err != nil {
		return err
	}

	// Generate a new salt
	newSalt, err := encryption.GenerateSalt("en_salt.txt")
	if err != nil {
		return err
	}

	// Hash the password with the new salt
	newPasswordHash, err := encryption.HashInput(password, []byte(newSalt))
	if err != nil {
		return err
	}

	// Hash the encoding with the new salt
	newEncodingHash, err := encryption.HashInput(normalizedEncoding, []byte(newSalt))
	if err != nil {
		return err
	}

	// Update the user's password hash and encoding hash
	if err := UpdateOwnUserHashes(newPasswordHash, newEncodingHash); err != nil {
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
	if err := encryption.EncryptFile("test.db", "en_test.db", ownUser.EncodingHash, salt); err != nil {
		return err
	}
	if err := encryption.EncryptFile("en_test.db", "en_test.db", ownUser.PasswordHash, salt); err != nil {
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

func GetContactByName(name string) (*entities.Contact, error) {
	contacts, err := prod_api.GetContacts(map[string]interface{}{"name": name})
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	if len(contacts) == 0 {
		return nil, nil
	}

	return &contacts[0], nil
}

func GetChatsAndInfo() ([]entities.ChatDTO, error) {
	chats, err := prod_api.GetChats(nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	var chatDTOsWithMessages []entities.ChatDTO
	for _, chat := range chats {
		contact, err := prod_api.GetContact(chat.ContactID)
		if err != nil {
			return nil, fmt.Errorf("failed to get contact: %w", err)
		}
		// Create chat DTO with messages
		chatDTO := entities.ChatDTO{
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

func GetChatWithContact(contactID string) (*entities.Chat, error) {
	var chats []entities.Chat
	var destContact *entities.Contact

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

	var chat entities.Chat
	if len(chats) > 1 { // Unexpected number of chats, this should not happen.
		return nil, fmt.Errorf("fatal: unexpected number of chats for single contact: %s", destContact.Name)
	} else if len(chats) == 0 { // Chat does not exist, create it.
		fmt.Printf("Creating new chat with %s.\n", destContact.Name)
		chat = entities.Chat{
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
