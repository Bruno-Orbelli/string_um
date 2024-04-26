package node

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"
	"time"

	"string_um/string/main/flags"
	"string_um/string/models"
	boots "string_um/string/networking/bootstrap"
	"string_um/string/networking/mdns"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"

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

func saveMessage(jsonString string) {
	reader := strings.NewReader(jsonString)
	for i := 0; i < 5; i++ {
		resp, err := http.Post("http://localhost:3000/messages/create", "application/json", reader)
		if err != nil || resp.StatusCode != 201 {
			fmt.Printf("Failed to save message, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		break
	}
}

func readMessage(r *bufio.Reader) {
	for i := 0; i < 5; i++ {
		str, err := r.ReadString('\n')
		fmt.Println(str)
		if err != nil {
			fmt.Printf("Failed to read message, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		if str == "" {
			return
		}
		if str != "\n" {
			saveMessage(str)
		}
		break
	}
}

func writeMessage(w *bufio.Writer, message models.Message) error {
	for i := 0; i < 5; i++ {
		if err := json.NewEncoder(w).Encode(message); err != nil {
			fmt.Printf("Failed to write message, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			continue
		}
		break
	}

	for i := 0; i < 5; i++ {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Failed to flush writer, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
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

		resp, err := http.Get("http://localhost:3000/messages?alreadySent=false")
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("Failed to get unsent messages, retrying in 15 seconds.")
			time.Sleep(time.Second * 15)
			continue
		}

		var unsentMessages []models.Message
		json.NewDecoder(resp.Body).Decode(&unsentMessages)

		for _, message := range unsentMessages {
			contactID := message.Chat.Contact.ID
			resp, err := http.Get(fmt.Sprintf("http://localhost:3000/contactAddresses?contactID=%s", contactID))
			if err != nil || resp.StatusCode != 200 {
				fmt.Println("Failed to get contact addresses, retrying in 15 seconds.")
				time.Sleep(time.Second * 15)
				continue // TODO: Try to resend the same message
			}

			var contactAddrs []models.ContactAddress
			json.NewDecoder(resp.Body).Decode(&contactAddrs)

			for _, contactAddr := range contactAddrs {
				potencialMa, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", contactAddr.ObservedAddress, contactAddr.Contact.ID))
				if err != nil {
					fmt.Println("Error when trying to convert multiaddress, getting another.")
					continue
				}
				if err := connectToPeer(ctx, host, potencialMa); err != nil {
					fmt.Printf("Error when connecting to peer: %s. Trying another multiaddress.\n", err)
					continue
				}

				contactID, err := peer.Decode(contactAddr.Contact.ID)

				if err != nil {
					fmt.Println("Couldn't decode ContactID, trying another multiaddress.")
					continue
				}
				s, err := openStreamToPeer(ctx, host, contactID)
				if err != nil {
					fmt.Printf("Couldn't open stream to peer: %s. Trying another multiaddress.\n", err)
					continue
				}
				if err := writeMessage(bufio.NewWriter(s), message); err != nil {
					fmt.Printf("Couldn't write message to peer: %s. Trying another multiaddress.\n", err)
					continue
				}

				// Mark message as sent
				for i := 0; i < 5; i++ {
					resp, err := http.Post(fmt.Sprintf("http://localhost:3000/messages/setAsSent/%s", message.ID), "application/json", nil)
					if err != nil || resp.StatusCode != 204 {
						fmt.Printf("Failed to mark message as sent, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
						time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
						continue
					}
					break
				}
			}
		}
	}
}

func connectToPeer(ctx context.Context, host host.Host, peerMa ma.Multiaddr) error {
	/* decodedPeerID, err := peer.Decode(peerID)
	if err != nil {
		return errors.New("couldn't decode PeerID")
	}

	for i := 0; i < 5; i++ {
		if len(host.Peerstore().PeerInfo(decodedPeerID).Addrs) == 0 {
			fmt.Printf("Failed to get an address at Peerstore, retrying after %f seconds.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			if i == 4 {
				continue
			} else {
				return errors.New("couldn't get an address for Peer")
			}
		}
		break
	}

	fmt.Println("Found address, attempting connection.") */

	addrInfo, err := peer.AddrInfoFromP2pAddr(peerMa)
	if err != nil {
		return errors.New("couldn't get AddrInfo from the provided multiaddress")
	}

	for i := 0; i < 5; i++ {
		if err := host.Connect(ctx, *addrInfo); err != nil {
			fmt.Printf("Failed to connect to peer, retrying after %f seconds.\n", math.Pow(float64(i+1), 2))
			time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
			if i == 4 {
				continue
			} else {
				return errors.New("couldn't connect to Peer")
			}
		}
		break
	}

	fmt.Printf("Connected to %s.\n", addrInfo.ID)
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

func AddNewContact(host host.Host, contactID string, contactName string) error {
	// Check if contact is in Peerstore
	contactPeerID, err := peer.Decode(contactID)
	if err != nil {
		return err
	}

	// Get known addresses for the contact
	var knownAddrs []ma.Multiaddr
	for i := 0; i < 6; i++ {
		knownAddrs = host.Peerstore().Addrs(contactPeerID)
		if len(knownAddrs) == 0 {
			fmt.Printf("Failed to get known addresses for contact, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			if i <= 4 {
				time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
				continue
			} else {
				return errors.New("couldn't get known addresses for contact")
			}
		}
		break
	}

	// Add the contact to the database
	newContact := models.Contact{
		ID:   contactPeerID.String(),
		Name: contactName,
	}
	newContactJSON, err := json.Marshal(newContact)
	if err != nil {
		return errors.New("couldn't marshal contact to JSON")
	}
	for i := 0; i < 5; i++ {
		resp, err := http.Post("http://localhost:3000/contacts/create", "application/json", bytes.NewReader(newContactJSON))
		if err != nil || resp.StatusCode != 201 {
			fmt.Printf("Failed to save contact, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
			if i > 4 {
				time.Sleep(time.Duration(math.Pow(float64(i+1), 2)) * time.Second)
				continue
			} else {
				return errors.New("couldn't save contact")
			}
		}
		break
	}

	// Add the contactAddrs to the database
	for _, addr := range knownAddrs {
		newAddr := models.ContactAddress{
			Contact:         newContact,
			ObservedAddress: addr.String(),
			ObservedAt:      time.Now(),
		}
		newAddrJSON, err := json.Marshal(newAddr)
		if err != nil {
			return errors.New("couldn't marshal contact address to JSON")
		}
		for i := 0; i < 5; i++ {
			resp, err := http.Post("http://localhost:3000/contactAddresses/create", "application/json", bytes.NewReader(newAddrJSON))
			if err != nil || resp.StatusCode != 201 {
				fmt.Println(resp.StatusCode)
				fmt.Printf("Failed to save contact address, retrying after %f second/s.\n", math.Pow(float64(i+1), 2))
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

func startLocalPeerDiscovery(host host.Host) {
	peerChan := mdns.InitMDNS(host)
	for { // Get all identified peers
		peer := <-peerChan
		host.Peerstore().AddAddrs(peer.ID, peer.Addrs, time.Minute*10)
		fmt.Println("New hosts received. Peers with addrs:", host.Peerstore().PeersWithAddrs())
	}
}

type myNotifiee struct {
	network.Notifiee
}

func (n *myNotifiee) Connected(_ network.Network, c network.Conn) {
	var nodeType string
	if slices.Contains(config.RelayAddresses, c.RemoteMultiaddr()) {
		nodeType = "relay"
	} else if slices.Contains(config.BootstrapPeers, c.RemoteMultiaddr()) {
		nodeType = "bootstrap"
	} else {
		nodeType = "peer"
	}
	fmt.Printf("Established connection with %s: %s\n", nodeType, c.RemotePeer())
}

func (n *myNotifiee) Disconnected(_ network.Network, c network.Conn) {
	var nodeType string
	if slices.Contains(config.RelayAddresses, c.RemoteMultiaddr()) {
		nodeType = "relay"
	} else if slices.Contains(config.BootstrapPeers, c.RemoteMultiaddr()) {
		nodeType = "bootstrap"
	} else {
		nodeType = "peer"
	}
	fmt.Printf("Connection with %s (%s) has been terminated.\n", nodeType, c.RemotePeer())
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
			fmt.Println("Failed to advertise, retrying in 15 seconds.")
			time.Sleep(time.Second * 15)
			continue
		}
		fmt.Println("Successfully advertised.")
		time.Sleep(time.Minute * 10)
	}
}
