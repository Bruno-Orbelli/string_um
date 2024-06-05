package main

import (
	"errors"
	"fmt"
	"os"
	"string_um/string/models"
	"time"

	prod_api "string_um/string/client/prod-api"
	"string_um/string/main/encryption"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"gorm.io/gorm"
)

func createOwnUser() (*models.OwnUser, *models.Contact, error) {
	fmt.Println("Creating own user...")

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

func getOwnUserIfExists() (*models.OwnUser, *models.Contact, error) {
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
		return createOwnUser()
	}
}

func addMessageToBeSent(chatID uuid.UUID, senderID string, body string) error {
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

func initDatabase(password string) error {
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
}

func getChatWithContact(contactID string) (*models.Chat, error) {
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

/*func main() {
	// Parse the command line arguments.
	config, err := flags.ParseFlags()
	if err != nil {
		panic(err)
	}

	// Handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithUseTransient(ctx, "relay info")

	// Run the application in a goroutine
	go func() {
		// Defer the cancellation of the context to ensure it is always called
		defer cancel()

		// Handle termination signal
		<-sigCh

		// Cancel the context
		cancel()
		<-ctx.Done()
	}()

	// Run the application
	if err := runApplication(ctx, config); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Encrypt the database file and remove the decrypted file.
	if err := encryption.EncryptFile("test.db", "en_test.db", config.Password, "en_salt.txt"); err != nil {
		panic(err)
	}
	if err := os.Remove("test.db"); err != nil {
		panic(err)
	}
}

func runApplication(ctx context.Context, config flags.Config) error {
	if err := InitDatabase(config.Password); err != nil {
		return err
	}

	// Run APIs according to config.
	go prod_api.RunDatabaseAPI()
	if config.Debug {
		go debug_api.RunDatabaseAPI()
	}
	time.Sleep(1 * time.Second)

	// Get own user, create if not exists.
	ownUser, ownUserContact, err := getOwnUserIfExists()
	if err != nil {
		return err
	}

	// Create a new host.
	fmt.Printf("Creating host with addresses: %s.\n", config.ListenAddresses)
	unmarshalledKey, err := crypto.UnmarshalPrivateKey(ownUser.PrivateKey)
	if err != nil {
		return err
	}
	host, _, err := node.CreateNewNode(
		ctx,
		unmarshalledKey,
		config.ListenAddresses,
		config.RelayAddresses,
		config.BootstrapPeers,
		protocol.ConvertFromStrings([]string{config.ProtocolID})[0],
	)
	if err != nil {
		return err
	}
	defer host.Close()

	// Set user ID and contact ID.
	hostinfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	fmt.Printf("Host info: %s.\n", hostinfo)

	// Send a message to the peer if a peer ID is provided.
	if config.PeerID != "" {
		if err := node.AddNewContact(host, config.PeerID, "juan23"); err != nil {
			fmt.Println("Failed to add contact:", err)
			return err
		}

		chat, err := getChatWithContact(config.PeerID)
		if err != nil {
			return err
		}

		// Add message to be sent.
		if err := addMessageToBeSent(chat.ID, ownUserContact.ID, "Hello, world!"); err != nil {
			return err
		}

		/* if config.PeerID != "" {
			fmt.Println("Trying to find peer...")
			decodedPeerID, err := peer.Decode(config.PeerID)
			if err != nil {
				panic(err)
			}
			for i := 0; i < 5; i++ {
				if len(host1.Peerstore().PeerInfo(decodedPeerID).Addrs) == 0 {
					fmt.Printf("Failed to get an address at Peerstore, retrying after %d seconds.\n", (i+1)*5)
					time.Sleep(time.Duration((i+1)*5) * time.Second)
					continue
				}
				fmt.Println("Found address, attempting connection.")
				if err = host1.Connect(ctx, host1.Peerstore().PeerInfo(decodedPeerID)); err != nil {
					panic(err)
				}
				break
			}

			if host1.Network().ConnsToPeer(decodedPeerID) != nil {
				// Open a stream with the given peer.
				s, err := host1.NewStream(ctx, decodedPeerID, "/chat/0.0.1")
				if err != nil {
					panic(err)
				}
				rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

				go writeData(rw)
				go readData(rw, decodedPeerID)
			}

		} else if len(config.RelayAddresses) > 0 {
			for _, relayAddr := range config.RelayAddresses {
				fmt.Println("Attempting connection to relay at", relayAddr)
				relayInfo, err := peer.AddrInfosFromP2pAddrs(relayAddr)
				if err != nil {
					panic(err)
				}
				if err = host1.Connect(ctx, relayInfo[0]); err != nil {
					fmt.Println("Failed to connect to relay, retrying with another.")
					time.Sleep(5 * time.Second)
					continue
				}
				_, err = client.Reserve(ctx, host1, relayInfo[0])
				if err != nil {
					fmt.Println("Failed to reserve relay, retrying with another.")
					time.Sleep(5 * time.Second)
					continue
				}
				fmt.Println("Awaiting incoming connections...")
			}
		} else {
			panic("no info to connect to a peer or relay")
		}
	}

	// Wait for cancellation signal or context done
	<-ctx.Done()

	return nil
}*/
