package node

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	prod_api "string_um/string/client/prod-api"
	"string_um/string/entities"
	"string_um/string/main/flags"
	"string_um/string/main/tui/globals"
	boots "string_um/string/networking/bootstrap"
	"string_um/string/networking/mdns"

	"github.com/google/uuid"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
	"gorm.io/gorm"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

func handleStream(stream network.Stream) {
	// Create a buffer stream for non-blocking read.
	reader := bufio.NewReader(stream)
	go readMessage(reader)
}

func unmarshallMessage(str string) (*entities.Message, error) {
	var message entities.Message
	err := json.Unmarshal([]byte(str), &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func alterAndSaveMessage(message entities.Message) {
	// Alter message to appear as sent
	message.AlreadySent = true
	if _, err := prod_api.CreateMessage(message); err != nil {
		// fmt.Printf("Failed to save message: %s.\n", err)
	}
}

func readMessage(r *bufio.Reader) {
	for i := 0; i < 5; i++ {
		str, err := r.ReadString('\n')
		if err != nil {
			// fmt.Printf("Failed to read message, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		if str == "" {
			return
		}
		if str != "\n" {
			message, err := unmarshallMessage(str)
			if err != nil {
				// fmt.Println("Failed to unmarshal message.")
				return
			}
			openNewChatIfNotExists(message.ChatID, message.SentByID)
			alterAndSaveMessage(*message)
		}
		globals.ChatsRefreshedChan <- true
		globals.MessagesRefreshedChan <- true
		break
	}
}

func writeMessage(w *bufio.Writer, message entities.Message) error {
	for i := 0; i < 5; i++ {
		if err := json.NewEncoder(w).Encode(message); err != nil {
			// fmt.Printf("Failed to write message, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 5; i++ {
		err := w.Flush()
		if err != nil {
			// fmt.Printf("Failed to flush writer, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		return nil
	}

	return errors.New("couldn't write message")
}

func checkUnsentMessagesAndSend(ctx context.Context, host host.Host) {
	for {
		time.Sleep(time.Second * 5)
		params := map[string]interface{}{"already_sent": "0"}
		unsentMessages, err := prod_api.GetMessages(params)
		if err != nil {
			// fmt.Printf("Failed to get unsent messages: %s\n. Retrying in 15 seconds.", err)
			time.Sleep(time.Second * 15)
			continue
		}

		for _, message := range unsentMessages {
			// Check if chat is open in node
			chat, err := prod_api.GetChat(message.ChatID)
			if err != nil {
				// fmt.Printf("Failed to get chat: %s\n. Retrying in 15 seconds.", err)
				time.Sleep(time.Second * 15)
				continue
			}

			// Get contact addresses for contact
			params = map[string]interface{}{"contact_id": chat.ContactID}
			contactAddrs, err := prod_api.GetContactAddresses(params)
			if err != nil {
				// fmt.Printf("Failed to get contactAddresses for %s: %s\n. Retrying in 15 seconds.", chat.ContactID, err)
				time.Sleep(time.Second * 15)
				continue // TODO: Try to resend the same message
			}

			for _, contactAddr := range contactAddrs {
				// Connect to peer
				potencialMa, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", contactAddr.ObservedAddress, contactAddr.ContactID))
				if err != nil {
					// fmt.Println("Error when trying to convert multiaddress, getting another.")
					continue
				}
				if err := connectToPeer(ctx, host, potencialMa); err != nil {
					// fmt.Printf("Error when connecting to peer: %s. Trying another multiaddress.\n", err)
					continue
				}

				// Open stream to peer
				contactID, err := peer.Decode(contactAddr.ContactID)
				if err != nil {
					//fmt.Println("Couldn't decode ContactID, trying another multiaddress.")
					continue
				}
				s, err := openStreamToPeer(ctx, host, contactID)
				if err != nil {
					// fmt.Printf("Couldn't open stream to peer: %s. Trying another multiaddress.\n", err)
					continue
				}

				// Write message to peer
				if err := writeMessage(bufio.NewWriter(s), message); err != nil {
					// fmt.Printf("Couldn't write message to peer: %s. Trying another multiaddress.\n", err)
					continue
				}

				// Mark message as sent
				params = map[string]interface{}{"already_sent": true}
				if _, err := prod_api.UpdateMessage(message.ID, params); err != nil {
					// fmt.Printf("Couldn't mark message %s as sent: %s. Trying another multiaddress.\n", message.ID, err)
					continue
				}

				break
			}
		}
	}
}

func connectToPeer(ctx context.Context, host host.Host, peerMa ma.Multiaddr) error {
	addrInfo, err := peer.AddrInfoFromP2pAddr(peerMa)
	if err != nil {
		return errors.New("couldn't get AddrInfo from the provided multiaddress")
	}

	for i := 0; i < 5; i++ {
		if err := host.Connect(ctx, *addrInfo); err != nil {
			// fmt.Printf("Failed to connect to peer, retrying after %f seconds.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			if i == 4 {
				continue
			} else {
				return errors.New("couldn't connect to Peer")
			}
		}
		break
	}

	// fmt.Printf("Connected to %s.\n", addrInfo.ID)
	return nil
}

func openStreamToPeer(ctx context.Context, host host.Host, peerId peer.ID) (*bufio.Writer, error) {
	if host.Network().ConnsToPeer(peerId) != nil {
		// Open a stream with the given peer.
		s, err := host.NewStream(ctx, peerId, "/chat/0.0.1")
		if err != nil {
			return nil, err
		}
		w := bufio.NewWriter(s)
		return w, nil
	} else {
		return nil, fmt.Errorf("no connection available to open stream to %s", peerId)
	}
}

var config = flags.Config{}

func CreateNewNode(ctx context.Context, priv crypto.PrivKey, listenAddrs []ma.Multiaddr, relayAddrs []ma.Multiaddr, bootsAddrs []ma.Multiaddr, protocolId protocol.ID) (host.Host, *dht.IpfsDHT, error) {
	// Set the configuration
	config.ListenAddresses = listenAddrs
	config.RelayAddresses = relayAddrs
	config.BootstrapPeers = bootsAddrs
	config.ProtocolID = protocol.ConvertToStrings([]protocol.ID{protocolId})[0]

	// Convert relay addresses to peer.AddrInfo
	relayAddrInfos := make([]peer.AddrInfo, len(relayAddrs))
	for i, relayAddr := range relayAddrs {
		addrInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
		if err != nil {
			return nil, nil, err
		}
		relayAddrInfos[i] = *addrInfo
	}

	// Convert bootstrap addresses to peer.AddrInfo
	bootsAddrInfos := make([]peer.AddrInfo, len(bootsAddrs))
	for i, bootsAddr := range bootsAddrs {
		addrInfo, err := peer.AddrInfoFromP2pAddr(bootsAddr)
		if err != nil {
			return nil, nil, err
		}
		bootsAddrInfos[i] = *addrInfo
	}

	var idht *dht.IpfsDHT

	// Create a new libp2p Host with Connection Manager
	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}

	node, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrs(listenAddrs...),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		libp2p.EnableAutoRelayWithStaticRelays(relayAddrInfos),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
	)
	if err != nil {
		return nil, nil, err
	}

	// Configure the stream handler to handle streams
	node.SetStreamHandler(protocolId, handleStream)
	node.Network().Notify(&myNotifiee{})

	// bootstrapHost(ctx, node, idht, bootsAddrInfos)

	go advertiseService(ctx, idht)
	go startLocalPeerDiscovery(node)
	go checkUnsentMessagesAndSend(ctx, node)

	return node, idht, nil
}

func AddKnownAddressesForContact(host host.Host, contactID string) error {
	// Check if the contactID is valid
	decodedContactID, err := peer.Decode(contactID)
	if err != nil {
		return err
	}

	// Get known addresses for the contact
	var knownAddrs []ma.Multiaddr
	for i := 0; i < 6; i++ {
		knownAddrs = host.Peerstore().Addrs(decodedContactID)
		if len(knownAddrs) == 0 {
			// fmt.Printf("Failed to get known addresses for contact, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			if i <= 4 {
				time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
				continue
			} else {
				return errors.New("couldn't get known addresses for contact")
			}
		}
		break
	}

	// Add the contactAddrs to the database if they don't already exist
	params := map[string]interface{}{"contact_id": contactID}
	existingContactAddrs, err := prod_api.GetContactAddresses(params)
	if err != nil {
		return err
	}

	var existingMaddrs []ma.Multiaddr
	for _, addr := range existingContactAddrs {
		existingMaddrs = append(existingMaddrs, ma.StringCast(addr.ObservedAddress))
	}

	for _, addr := range knownAddrs {
		newAddr := entities.ContactAddress{
			ContactID:       contactID,
			ObservedAddress: addr.String(),
			ObservedAt:      time.Now(),
		}
		for i := 0; i < 5; i++ {
			_, err := prod_api.CreateContactAddress(newAddr)
			if err != nil {
				// fmt.Printf("Failed to save contact address: %s. Retrying after %f second/s.\n", createdAddr, math.Pow(float64(i+1), 2))
				if i > 4 {
					time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
					continue
				} else {
					return errors.New("couldn't save contact address")
				}
			}
			break
		}
	}

	return nil
}

func AddNewContact(host host.Host, contactID string, contactName string) error {
	// Check if the contactID is valid
	decodedContactID, err := peer.Decode(contactID)
	if err != nil {
		return err
	}

	// Check if the contact already exists or has default name
	for i := 0; i < 6; i++ {
		contact, err := prod_api.GetContact(contactID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// fmt.Printf("Failed to get contact: %s. Retrying after %f second/s.\n", err, math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			if i < 5 {
				continue
			} else {
				return errors.New("couldn't get contact")
			}
		}
		// Conditions to add or update contact
		if contact == nil { // Add the contact to the database
			newContact := entities.Contact{
				ID:   decodedContactID.String(),
				Name: contactName,
			}
			for i := 0; i < 5; i++ {
				_, err := prod_api.CreateContact(newContact)
				if err != nil {
					// fmt.Printf("Failed to save contact, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
					if i > 4 {
						time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
						continue
					} else {
						return errors.New("couldn't save contact")
					}
				}
				break
			}
		} else if contact.Name == contactID { // Contact already exists with default name
			params := map[string]interface{}{"name": contactName}
			for i := 0; i < 6; i++ {
				if _, err := prod_api.UpdateContact(contactID, params); err != nil {
					// fmt.Printf("Failed to update contact: %s. Retrying after %f second/s.\n", err, math.Pow(float64(i+1), 2))
					time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
					if i < 5 {
						continue
					} else {
						return errors.New("couldn't update contact")
					}
				}
			}
		}
		err = AddKnownAddressesForContact(host, contactID)
		return err // Contact now exists with custom name or has been updated
	}
	return nil // Contact now exists with custom name or has been updated
}

func SaveContactIfNotExists(contactID string) (*entities.Contact, error) {
	// TODO: Better error handling
	// Check if contact exists
	contact, err := prod_api.GetContact(contactID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if contact != nil {
		return contact, nil
	} else { // Contact doesn't exist, create it
		contact := entities.Contact{
			ID:   contactID,
			Name: contactID,
		}
		for i := 0; i < 6; i++ {
			createdContact, err := prod_api.CreateContact(contact)
			if err != nil {
				// fmt.Printf("Failed to create contact: %s. Retrying in 15 seconds.\n", err)
				time.Sleep(time.Second * 15)
				if i < 5 {
					continue
				} else {
					return nil, errors.New("couldn't create contact")
				}
			}
			return createdContact, nil
		}
	}
	return nil, errors.New("couldn't create contact")
}

func openNewChatIfNotExists(chatID uuid.UUID, contactID string) error {
	contact, err := SaveContactIfNotExists(contactID)
	if err != nil {
		return err
	}

	// Check if chat already exists
	var chat *entities.Chat
	for i := 0; i < 6; i++ {
		chat, err = prod_api.GetChat(chatID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// fmt.Printf("Failed to get chat: %s, retrying in 15 seconds.\n", err)
			time.Sleep(time.Second * 15)
			if i < 5 {
				continue
			} else {
				return errors.New("couldn't get chats")
			}
		}
		break
	}

	if chat != nil && chat.ContactID != contactID { // Chat already exists but is associated with another contact
		return errors.New("chat mismatch: chat already exists but is associated with another contact") // In theory, this should never happen due to UUIDs being unique
	} else if chat == nil { // Chat doesn't exist, create it
		// fmt.Printf("Creating new chat with %s.\n", contact.Name)
		chat := entities.Chat{
			ID:        chatID,
			ContactID: contact.ID,
		}
		for i := 0; i < 6; i++ {
			_, err = prod_api.CreateChat(chat)
			if err != nil {
				// fmt.Printf("Failed to create chat: %s. Retrying in 15 seconds.\n", err)
				time.Sleep(time.Second * 15)
				if i < 5 {
					continue
				} else {
					return errors.New("couldn't create chat")
				}
			}
			return nil
		}
	}
	return nil // Chat already exists and is associated with the correct contact
}

func startLocalPeerDiscovery(host host.Host) {
	peerChan := mdns.InitMDNS(host)
	for { // Get all identified peers
		peer := <-peerChan
		host.Peerstore().AddAddrs(peer.ID, peer.Addrs, time.Minute*10)
		// fmt.Println("New hosts received. Peers with addrs:", host.Peerstore().PeersWithAddrs())
	}
}

type myNotifiee struct {
	network.Notifiee
}

func (n *myNotifiee) Connected(_ network.Network, c network.Conn) {
	var _ string
	/*if slices.Contains(config.RelayAddresses, c.RemoteMultiaddr()) {
		nodeType = "relay"
	} else if slices.Contains(config.BootstrapPeers, c.RemoteMultiaddr()) {
		nodeType = "bootstrap"
	} else {
		nodeType = "peer"
	}*/
	// fmt.Printf("Established connection with %s: %s\n", nodeType, c.RemotePeer())
}

func (n *myNotifiee) Disconnected(_ network.Network, c network.Conn) {
	var _ string
	/*if slices.Contains(config.RelayAddresses, c.RemoteMultiaddr()) {
		nodeType = "relay"
	} else if slices.Contains(config.BootstrapPeers, c.RemoteMultiaddr()) {
		nodeType = "bootstrap"
	} else {
		nodeType = "peer"
	}*/
	// fmt.Printf("Connection with %s (%s) has been terminated.\n", nodeType, c.RemotePeer())
}

func (n *myNotifiee) ListenClose(_ network.Network, mAddr ma.Multiaddr) {
	// fmt.Printf("Clossing listener at multiaddress %s.\n", mAddr)
}

func bootstrapHost(ctx context.Context, host host.Host, dhtInstance *dht.IpfsDHT, bootstrapPeers []peer.AddrInfo) error {
	if err := boots.TestBootstraps(ctx, host, bootstrapPeers); err != nil {
		return err
	}
	dhtInstance.Bootstrap(ctx)
	fmt.Println("Succesfully bootstraped to the network.")

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := host.Addrs()
	fmt.Println("This host is reachable at the following bootstraped addresses: ")
	for _, addr := range addrs {
		fmt.Println(addr.Encapsulate(hostAddr))
	}

	return nil
}

func advertiseService(ctx context.Context, idht *dht.IpfsDHT) {
	for {
		announcement := discovery.NewRoutingDiscovery(idht)
		_, err := announcement.Advertise(ctx, "example-discovery")
		if err != nil {
			// fmt.Println("Failed to advertise, retrying in 15 seconds.")
			time.Sleep(time.Second * 15)
			continue
		}
		// fmt.Println("Successfully advertised.")
		time.Sleep(time.Minute * 10)
	}
}
