package main

import (
	"context"
	"flag"
	"fmt"
	mrand "math/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	dual "github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	help := flag.Bool("help", false, "Display Help")
	listenHost := flag.String("host", "0.0.0.0", "The bootstrap node host listen address\n")
	port := flag.Int("port", 4001, "The bootstrap node listen port")
	flag.Parse()

	if *help {
		fmt.Printf("This is a simple bootstrap node for kad-dht application using libp2p\n\n")
		fmt.Printf("Usage: \n   Run './bootnode'\nor Run './bootnode -host [host] -port [port]'\n")

		os.Exit(0)
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", *listenHost, *port)

	ctx := context.Background()
	r := mrand.New(mrand.NewSource(int64(*port)))

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", *listenHost, *port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	// dstore := dsync.MutexWrap(ds.NewMapDatastore())

	kad, err := dual.New(ctx, host)
	if err != nil {
		panic(err)
	}

	rhost := rhost.Wrap(host, kad)
	fmt.Println(rhost.ID(), rhost.Addrs())

	// Register a handler for connection events
	rhost.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
			fmt.Println("New connection established with peer:", c.RemotePeer())
			kad.LAN.PutValue(ctx, c.RemotePeer().String(), c.RemoteMultiaddr().Bytes())
			fmt.Println(kad.LAN.GetValue(ctx, c.RemotePeer().String()))
			// Print the current peer list
			fmt.Println("Current peer list:")
			for _, p := range kad.LAN.RoutingTable().ListPeers() {
				fmt.Println(p)
			}
			fmt.Println("")
		},
	})

	fmt.Println("")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/ipfs/%s\n", *listenHost, *port, rhost.ID())
	fmt.Println("")
	select {}
}
