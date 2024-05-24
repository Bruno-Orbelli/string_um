package main

import (
	"context"
	"errors"
	"fmt"
	"main/flags"
	"os"
	"os/signal"
	"string_um/string/models"
	"string_um/string/networking/node"
	"syscall"
	"time"

	debug_api "string_um/string/client/debug-api"
	prod_api "string_um/string/client/prod-api"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"gorm.io/gorm"
)

func createOwnUser(password string) (*models.OwnUser, *models.Contact, error) {
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
		Password:   password,
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

func getOwnUserIfExists(config flags.Config) (*models.OwnUser, *models.Contact, error) {
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
		return createOwnUser(config.Password)
	}
}

func main() {
	// Parse the command line arguments.
	config, err := flags.ParseFlags()
	if err != nil {
		panic(err)
	}

	// Run APIs according to config.
	go prod_api.RunDatabaseAPI()
	if config.Debug {
		go debug_api.RunDatabaseAPI()
	}
	time.Sleep(1 * time.Second)

	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithUseTransient(ctx, "relay info")
	defer cancel()

	// Get own user, create if not exists.
	ownUser, ownUserContact, err := getOwnUserIfExists(config)
	if err != nil {
		panic(err)
	}

	// Create a new host.
	fmt.Printf("Creating host with addresses: %s.\n", config.ListenAddresses)
	unmarshalledKey, err := crypto.UnmarshalPrivateKey(ownUser.PrivateKey)
	if err != nil {
		panic(err)
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
		panic(err)
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
			panic(err)
		}

		var destContact *models.Contact
		var chats []models.Chat

		// Get contact to send message to.
		destContact, err = prod_api.GetContact(config.PeerID)
		if err != nil {
			panic("failed to get contact: " + err.Error())
		}

		// Get chat with contact.
		params := map[string]interface{}{"contact_id": destContact.ID}
		chats, err = prod_api.GetChats(params)
		if err != nil {
			panic("failed to get chats: " + err.Error())
		}

		var chat models.Chat
		if len(chats) > 1 { // Unexpected number of chats, this should not happen.
			panic(errors.New("fatal: unexpected number of chats for single contact: " + destContact.Name))
		} else if len(chats) == 0 { // Chat does not exist, create it.
			fmt.Printf("Creating new chat with %s.\n", destContact.Name)
			chat = models.Chat{
				ContactID: destContact.ID,
			}
			createdChat, err := prod_api.CreateChat(chat)
			if err != nil {
				panic("failed to create chat: " + err.Error())
			}
			chat = *createdChat
		} else { // Chat exists, use it.
			chat = chats[0]
		}

		message := models.Message{
			ID:          uuid.New(),
			ChatID:      chat.ID,
			AlreadySent: false,
			SentByID:    ownUserContact.ID,
			SentAt:      time.Now(),
			Message:     "Hello, world!",
		}
		createdMessage, err := prod_api.CreateMessage(message)
		if err != nil {
			panic("failed to create message: " + err.Error())
		}

		fmt.Printf("Created message with ID=%s.\n", createdMessage.ID)

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
		} */
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}
