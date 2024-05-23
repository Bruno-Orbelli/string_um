package flags

import (
	"flag"
	"strings"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"
)

type addrList []maddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

type Config struct {
	Password        string
	ListenAddresses addrList
	BootstrapPeers  addrList
	RelayAddresses  addrList
	PeerID          string
	ProtocolID      string
}

func ParseFlags() (Config, error) {
	config := Config{}
	flag.StringVar(&config.Password, "password", "", "Sets a password for the user")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.Var(&config.BootstrapPeers, "boots", "Adds a peer multiaddress to the bootstrap list")
	flag.Var(&config.RelayAddresses, "relay", "Adds a relay multiaddress to the relay list")
	flag.StringVar(&config.PeerID, "peer", "", "Sets a peer id to connect to")
	flag.StringVar(&config.ProtocolID, "pid", "/chat/0.0.1", "Sets a protocol id for stream headers")
	flag.Parse()

	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	return config, nil
}