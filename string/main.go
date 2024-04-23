package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

func main() {
	// Parse the command line arguments.
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithUseTransient(ctx, "relay info")
	defer cancel()

	// Create a new host.
	fmt.Printf("Creating host with addresses: %s.\n", config.ListenAddresses)
	host1, _, err := CreateNewNode(ctx, config.ListenAddresses, config.RelayAddresses, config.BootstrapPeers, protocol.ConvertFromStrings([]string{config.ProtocolID})[0])
	if err != nil {
		panic(err)
	}
	defer host1.Close()

	hostinfo := peer.AddrInfo{
		ID:    host1.ID(),
		Addrs: host1.Addrs(),
	}

	fmt.Printf("Host info: %s.\n", hostinfo)
	host1.SetStreamHandler("/chat/0.0.1", handleStream)

	if config.PeerID != "" {
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}
