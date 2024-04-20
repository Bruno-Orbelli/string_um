package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"

	ma "github.com/multiformats/go-multiaddr"
)

var H1_ADDRESSES = []string{
	"/ip4/192.168.54.225/tcp/50004",
	"/ip4/127.0.0.1/tcp/50000",
	"/ip6/::1/tcp/50001",
	"/ip4/127.0.0.1/udp/50002/quic",
	"/ip6/::1/udp/50003/quic",
}

var H2_ADDRESSES = []string{
	"/ip4/127.0.0.1/tcp/50001",
	"/ip6/::1/tcp/50006",
	"/ip4/127.0.0.1/udp/50007/quic",
	"/ip6/::1/udp/50008/quic",
	"/ip4/192.168.54.225/tcp/50009",
}

var RELAY_ADDRESSES = []string{
	"/ip4/192.168.54.225/tcp/51001",
}

var PUBLIC_RELAY_ADDRESSES = []string{
	"/ip4/190.15.209.53/tcp/51001",
}

var PUBLIC_BOOTSTRAP_ADDRESSES = []string{
	"/ip4/190.15.209.53/tcp/51000",
}

func main() {
	// log.SetAllLoggers(log.LevelDebug)

	ownAddr := flag.String("own-address", H1_ADDRESSES[0], "own address")
	peerAddr := flag.String("peer-address", "", "peer address")
	peerId := flag.String("peer-id", "", "peer id")
	relayAddr := flag.String("relay-address", PUBLIC_RELAY_ADDRESSES[0], "relay address")
	relayId := flag.String("relay-id", "", "relay id")
	bootstrapId := flag.String("bootstrap-id", "", "bootstrap id")
	flag.Parse()

	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	ctx = network.WithUseTransient(ctx, "relay info")
	defer cancel()

	// Build the list of addresses
	ownAddrList := []string{*ownAddr}

	// Build the relay address
	relayMa := ma.StringCast(fmt.Sprintf("%s/p2p/%s", *relayAddr, *relayId))
	relayInfo, err := peer.AddrInfoFromP2pAddr(relayMa)
	if err != nil {
		panic(err)
	}

	// Build the bootstrap addresses
	bootstrapInfos := make([]peer.AddrInfo, len(PUBLIC_BOOTSTRAP_ADDRESSES))
	for i, addr := range PUBLIC_BOOTSTRAP_ADDRESSES {
		bootstrapMa := ma.StringCast(fmt.Sprintf("%s/p2p/%s", addr, *bootstrapId))
		bootstrapInfo, err := peer.AddrInfoFromP2pAddr(bootstrapMa)
		if err != nil {
			panic(err)
		}
		bootstrapInfos[i] = *bootstrapInfo
	}

	// Create a new host.
	fmt.Printf("Creating host with addresses: %s.\n", *ownAddr)
	host1, _, err := CreateNewNode(ctx, ownAddrList, *relayInfo, bootstrapInfos)
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

	if *peerAddr != "" {
		fmt.Println("Dialing to", *peerAddr)
		peerMA, err := ma.NewMultiaddr(*peerAddr)
		if err != nil {
			panic(err)
		}
		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			panic(err)
		}

		// Connect to the node at the given address.
		if err := host1.Connect(ctx, *peerAddrInfo); err != nil {
			fmt.Println(err.Error())
			host1.Connect(ctx, *relayInfo)
			newRemoteMa, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s/p2p-circuit/p2p/%s", peerAddrInfo.Addrs[0], relayInfo.ID, peerAddrInfo.ID))
			if err != nil {
				panic(err)
			}

			newPeerAddrInfo, err := peer.AddrInfoFromP2pAddr(newRemoteMa)
			if err != nil {
				panic(err)
			}
			// peerAddrInfo.Addrs = append(peerAddrInfo.Addrs, newRemoteMa)
			_, err = client.Reserve(ctx, host1, *relayInfo)
			if err != nil {
				panic(err)
			}
			if err := host1.Connect(ctx, *newPeerAddrInfo); err != nil {
				panic(err)
			}
		}

		// Open a stream with the given peer.
		s, err := host1.NewStream(ctx, peerAddrInfo.ID, "/chat/0.0.1")
		if err != nil {
			panic(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go writeData(rw)
		go readData(rw, peerAddrInfo.ID)

	} else if *peerId != "" {
		decodedPeerId, err := peer.Decode(*peerId)
		if err != nil {
			panic(err)
		}
		/* peerAddrInfo, err := dht1.FindPeer(ctx, decodedPeerId)
		if err != nil {
			panic(err)
		}
		if len(peerAddrInfo.Addrs) != 0 {
			fmt.Println("Found peer at", peerAddrInfo.Addrs[0])
			if err = host1.Connect(ctx, peerAddrInfo); err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Failed to find peer at", decodedPeerId)
		} */
		for i := 0; i < 3; i++ {
			fmt.Println("Trying to find peer...")
			if len(host1.Peerstore().PeerInfo(decodedPeerId).Addrs) == 0 {
				fmt.Printf("Failed to get an address at Peerstore, retrying after %d seconds.\n", (i+1)*5)
				time.Sleep(time.Duration((i+1)*5) * time.Second)
				continue
			}
			fmt.Println("Found address, attempting connection.")
			if err = host1.Connect(ctx, host1.Peerstore().PeerInfo(decodedPeerId)); err != nil {
				panic(err)
			}
			break
		}
		// Open a stream with the given peer.
		s, err := host1.NewStream(ctx, decodedPeerId, "/chat/0.0.1")
		if err != nil {
			panic(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go writeData(rw)
		go readData(rw, decodedPeerId)

	} else {
		fmt.Println("Connecting to relay at", relayInfo.String())
		if err = host1.Connect(ctx, *relayInfo); err != nil {
			panic(err)
		}
		_, err = client.Reserve(ctx, host1, *relayInfo)
		if err != nil {
			panic(err)
		}
		fmt.Println("Awaiting incoming connections...")
	}

	// TODO: Set the hosts on different goroutines to listen for incoming connections.
	// TODO: Relay
	// TODO: kdm

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}
