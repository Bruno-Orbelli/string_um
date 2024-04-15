package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "/example/1.0.0"
const discoveryNamespace = "example"

func main() {
	// Add -peer-address flag
	peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()

	// Create the libp2p host.
	//
	// Note that we are explicitly passing the listen address and restricting it to IPv4 over the
	// loopback interface (127.0.0.1).
	//
	// Setting the TCP port as 0 makes libp2p choose an available port for us.
	// You could, of course, specify one if you like.
	multi, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/50000")
	host, err := libp2p.New(libp2p.ListenAddrs(multi))
	if err != nil {
		panic(err)
	}
	defer host.Close()

	// Print this node's full address
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	addr, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Node addresses: ", addr)

	// Setup peer discovery.
	discoveryService := mdns.NewMdnsService(
		host,
		discoveryNamespace,
		&discoveryNotifee{h: host},
	)
	discoveryService.Start()
	defer discoveryService.Close()

	// Setup AutoNAT service.
	autoNatService := setupAutoNAT(host)
	fmt.Println(autoNatService.Status())
	defer autoNatService.Close()

	// Setup a stream handler.
	//
	// This gets called every time a peer connects and opens a stream to this node.
	host.SetStreamHandler(protocolID, func(s network.Stream) {
		go writeMessage(s)
		go readMessage(s)
	})

	// If we received a peer address, we should connect to it.
	if *peerAddr != "" {
		// Parse the multiaddr string.
		peerMA, err := multiaddr.NewMultiaddr(*peerAddr)
		if err != nil {
			panic(err)
		}
		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			panic(err)
		}

		// Connect to the node at the given address.
		if err := host.Connect(context.Background(), *peerAddrInfo); err != nil {
			panic(err)
		}
		fmt.Println("Connected to", peerAddrInfo.String())

		// Open a stream with the given peer.
		s, err := host.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
		if err != nil {
			panic(err)
		}

		// Start the write and read threads.
		go writeMessage(s)
		go readMessage(s)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}

func setupAutoNAT(host host.Host) autonat.AutoNAT {
	// Setup AutoNAT service.
	//
	// This will attempt to automatically open ports on the NAT.
	autoNatService, err := autonat.New(host)
	if err != nil {
		panic(err)
	}
	return autoNatService
}

func writeMessage(s network.Stream) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input:", err)
			return
		}

		// Remove newline character from input
		input = strings.TrimSpace(input)

		// Write message length as uint16 followed by message
		msgLen := uint16(len(input))
		if err := binary.Write(s, binary.BigEndian, msgLen); err != nil {
			fmt.Println("An error occurred while writing message length:", err)
			return
		}

		if _, err := s.Write([]byte(input)); err != nil {
			fmt.Println("An error occurred while writing message:", err)
			return
		}
	}
}

func readMessage(s network.Stream) {
	reader := bufio.NewReader(s)
	for {
		var msgLen uint16
		if err := binary.Read(s, binary.BigEndian, &msgLen); err != nil {
			fmt.Println("An error occurred while reading message length:", err)
			return
		}

		msg := make([]byte, msgLen)
		if _, err := reader.Read(msg); err != nil {
			fmt.Println("An error occurred while reading message:", err)
			return
		}

		fmt.Printf("\n> %s: %s\n", s.Conn().RemotePeer(), string(msg))
	}
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	fmt.Println("found peer", peerInfo.String())
}
