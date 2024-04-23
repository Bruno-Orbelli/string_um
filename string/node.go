package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

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
	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go readData(rw, stream.Conn().RemotePeer())
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter, remotePeer peer.ID) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s: %s\x1b[0m> ", remotePeer, str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

var config = Config{}

func CreateNewNode(ctx context.Context, listenAddrs []ma.Multiaddr, relayAddrs []ma.Multiaddr, bootsAddrs []ma.Multiaddr, protocolId protocol.ID) (host.Host, *dht.IpfsDHT, error) {
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

	// Generate a key pair for the host
	priv, _, err := crypto.GenerateKeyPair(
		crypto.RSA,
		2048,
	)
	if err != nil {
		return nil, nil, err
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

	bootstrapHost(ctx, node, idht, bootsAddrInfos)

	go advertiseService(ctx, idht)
	go startLocalPeerDiscovery(node)

	return node, idht, nil
}

func startLocalPeerDiscovery(host host.Host) {
	peerChan := initMDNS(host)
	for { // Get all identified peers
		peer := <-peerChan
		if len(host.Peerstore().PeerInfo(peer.ID).Addrs) == 0 {
			host.Peerstore().AddAddrs(peer.ID, peer.Addrs, time.Minute*10)
		}
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
	fmt.Printf("Connection with %s (%s) has been terminated.", nodeType, c.RemotePeer())
}

func bootstrapHost(ctx context.Context, host host.Host, dhtInstance *dht.IpfsDHT, bootstrapPeers []peer.AddrInfo) error {
	if err := testBootstraps(ctx, host, bootstrapPeers); err != nil {
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
