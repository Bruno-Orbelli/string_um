package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/flags"
	"net/http"
	"os"
	"os/signal"
	"string_um/string/models"
	"string_um/string/networking/node"
	"syscall"
	"time"

	dapi "string_um/string/client/debug-api"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func createOwnUserIfNotExists(config flags.Config) (*models.OwnUser, *models.Contact, error) {
	// Check if own user exists
	resp, err := http.Get("http://localhost:3000/ownUser")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Own user exists, decode response
		var ownUser models.OwnUser
		if err := json.NewDecoder(resp.Body).Decode(&ownUser); err != nil {
			return nil, nil, err
		}

		// Retrieve own user contact
		resp, err = http.Get("http://localhost:3000/contacts?name=Me")
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		var ownUserContact models.Contact
		if err := json.NewDecoder(resp.Body).Decode(&ownUserContact); err != nil {
			return nil, nil, err
		}

		return &ownUser, &ownUserContact, nil

	case http.StatusNotFound:
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
			Password:   config.Password,
			PrivateKey: marshalledKey,
		}
		ownUserJSON, err := json.Marshal(ownUser)
		if err != nil {
			return nil, nil, err
		}

		// POST own user
		resp, err = http.Post("http://localhost:3000/ownUser/create", "application/json", bytes.NewReader(ownUserJSON))
		if err != nil || resp.StatusCode != http.StatusCreated {
			return nil, nil, errors.New("failed to create own user")
		}

		// Create own user contact
		ownUserContact := models.Contact{
			ID:   ownUser.ID,
			Name: "Me",
		}
		ownUserContactJSON, err := json.Marshal(ownUserContact)
		if err != nil {
			return nil, nil, err
		}

		// POST own user contact
		resp, err = http.Post("http://localhost:3000/contacts/create", "application/json", bytes.NewReader(ownUserContactJSON))
		if err != nil || resp.StatusCode != http.StatusCreated {
			return nil, nil, errors.New("failed to create own user contact")
		}

		return &ownUser, &ownUserContact, nil

	default:
		return nil, nil, errors.New("unexpected status code: " + resp.Status)
	}
}

func main() {
	// Create a new API server.
	go dapi.RunDatabaseAPI()
	time.Sleep(1 * time.Second)

	// Parse the command line arguments.
	config, err := flags.ParseFlags()
	if err != nil {
		panic(err)
	}

	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithUseTransient(ctx, "relay info")
	defer cancel()

	// Create own user if not exists.
	ownUser, ownUserContact, err := createOwnUserIfNotExists(config)
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

		var destContact models.Contact
		var chat models.Chat

		// Get contact to send message to.
		resp, err := http.Get(fmt.Sprintf("http://localhost:3000/contacts/%s", config.PeerID))
		if err != nil || resp.StatusCode != http.StatusOK {
			panic(errors.New("unexpected error or status code: " + resp.Status))
		}
		if err = json.NewDecoder(resp.Body).Decode(&destContact); err != nil {
			panic(err)
		}

		// Get chat with contact.
		resp, err = http.Get(fmt.Sprintf("http://localhost:3000/chats?contact_id=%s", destContact.ID))
		if err != nil || resp.StatusCode != http.StatusOK {
			panic(errors.New("unexpected error or status code: " + resp.Status))
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		if string(respBytes) == "[]\n" {
			fmt.Printf("Creating new chat with %s.\n", destContact.Name)
			chat = models.Chat{
				ContactID: destContact.ID,
			}
			chatJSON, err := json.Marshal(chat)
			if err != nil {
				panic(err)
			}
			resp, err = http.Post("http://localhost:3000/chats/create", "application/json", bytes.NewReader(chatJSON))
			if err != nil || resp.StatusCode != http.StatusCreated {
				panic(err)
			}
			respBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
		} else if resp.StatusCode != http.StatusOK {
			panic(errors.New("unexpected status code: " + resp.Status))
		}
		if err = json.NewDecoder(bytes.NewReader(respBytes)).Decode(&chat); err != nil {
			panic(err)
		}

		message := models.Message{
			ID:          uuid.New(),
			ChatID:      chat.ID,
			AlreadySent: false,
			SentByID:    ownUserContact.ID,
			SentAt:      time.Now(),
			Message:     "Hello, world!",
		}
		messageJSON, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		resp, err = http.Post("http://localhost:3000/messages/create", "application/json", bytes.NewReader(messageJSON))
		if err != nil || resp.StatusCode != http.StatusCreated {
			panic(err)
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
		} */
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}
