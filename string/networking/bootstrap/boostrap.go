package bootstrap

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multiaddr"
)

func RunBootstrap() {
	help := flag.Bool("help", false, "Display Help")
	listenHost := flag.String("host", "0.0.0.0", "The bootstrap node host listen address\n")
	port := flag.Int("port", 4001, "The bootstrap node listen port")
	flag.Parse()

	if *help {
		fmt.Printf("This is a simple bootstrap node for kad-dht application using libp2p\n\n")
		fmt.Printf("Usage: \n   Run './bootstrap'\nor Run './boostrap -host [host] -port [port]'\n")

		os.Exit(0)
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", *listenHost, *port)

	ctx := context.Background()

	// Creates a new RSA key pair for this host.
	prv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	var ip4Or6 string
	if strings.Contains(*listenHost, ".") {
		ip4Or6 = "ip4"
	} else if strings.Contains(*listenHost, ":") {
		ip4Or6 = "ip6"
	} else {
		panic("Invalid IP address.")
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/%s/%s/tcp/%d", ip4Or6, *listenHost, *port))
	var idht *dht.IpfsDHT

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prv),
		libp2p.EnableNATService(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		libp2p.DefaultTransports,
	)

	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", *listenHost, *port, host.ID())
	fmt.Println("")
	select {}
}
