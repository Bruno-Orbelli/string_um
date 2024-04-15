package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"

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

	go readData(rw, stream.Conn().RemotePeer())
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
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

func CreateNewNode(ctx context.Context, strAddrs []string, relayAddrInfo peer.AddrInfo) (host.Host, error) {
	relayAddrArray := [1]peer.AddrInfo{
		relayAddrInfo,
	}

	priv, _, err := crypto.GenerateKeyPair(
		crypto.RSA, // Select your key type. Ed25519 are nice short
		2048,       // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		return nil, err
	}

	var idht *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	node, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(strAddrs...),
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
		libp2p.EnableAutoRelayWithStaticRelays(relayAddrArray[:]),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
	)
	if err != nil {
		return nil, err
	}

	// Configure the stream handler to handle streams
	node.SetStreamHandler("/chat/0.0.1", handleStream)
	node.Network().Notify(&myNotifiee{})
	return node, nil
}

type myNotifiee struct {
	network.Notifiee
}

func (n *myNotifiee) Connected(_ network.Network, c network.Conn) {
	if !slices.Contains(PUBLIC_RELAY_ADDRESSES, c.RemoteMultiaddr().String()) {
		fmt.Println("Established connection with node:", c.RemotePeer())
	} else {
		fmt.Println("Established connection with relay:", c.RemotePeer())
	}
}

func (n *myNotifiee) Disconnected(_ network.Network, c network.Conn) {
	if !slices.Contains(PUBLIC_RELAY_ADDRESSES, c.RemoteMultiaddr().String()) {
		fmt.Println("Connection with", c.RemotePeer(), "has been terminated.")
	} else {
		fmt.Println("Connection with relay", c.RemotePeer(), "has been terminated.")
	}
}
