package main

import (
	"context"
	"fmt"
	"log"

	// golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"

	"github.com/libp2p/go-libp2p/core/host"

	"github.com/libp2p/go-libp2p/core/crypto"
)

var RELAY_ADDRESSES = []string{
	"/ip6/::/tcp/51001",
}

var PUBLIC_RELAY_ADDRESSES = []string{
	"/ip6/2002:c833:29b0:2:f816:3eff:fed5:bda/tcp/51001",
}

func main() {
	// golog.SetAllLoggers(golog.LevelDebug)

	fmt.Printf("Creating relay node with private addresses: %s and public addresses %s.\n", RELAY_ADDRESSES, PUBLIC_RELAY_ADDRESSES)
	relay, err := CreateRelayWithAutoNAT(RELAY_ADDRESSES)
	if err != nil {
		log.Fatal(err)
	}
	defer relay.Close()

	relayInfo := peer.AddrInfo{
		ID:    relay.ID(),
		Addrs: relay.Addrs(),
	}

	fmt.Printf("Relay info: %s.\n", relayInfo)

	select {}
}

func CreateRelayWithAutoNAT(strAddrs []string) (host.Host, error) {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new identity for the relay host
	priv, _, err := crypto.GenerateKeyPair(
		crypto.RSA, // Select your key type. Ed25519 are nice short
		2048,       // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		return nil, err
	}

	var idht *dht.IpfsDHT

	// Create a new libp2p Host that listens the specified addresses and enable the relay service.
	relayHost, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(strAddrs...),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		libp2p.EnableRelayService(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		// libp2p.ForceReachabilityPublic(),
	)
	if err != nil {
		return nil, err
	}

	_, err = relay.New(relayHost)
	if err != nil {
		return nil, err
	}
	_, err = autonat.New(relayHost)
	if err != nil {
		return nil, err
	}

	relayHost.Network().Notify(&myNotifiee{})
	relayHost.ConnManager().Notifee()

	return relayHost, nil
}

type myNotifiee struct {
	network.Notifiee
}

func (n *myNotifiee) Connected(_ network.Network, c network.Conn) {
	fmt.Println("A new node has connected:", c.RemotePeer())
}

func (n *myNotifiee) Disconnected(_ network.Network, c network.Conn) {
	fmt.Println(c.RemotePeer(), "has disconnected.")
}
