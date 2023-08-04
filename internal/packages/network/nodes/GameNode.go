package nodes

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"

	// ma "github.com/multiformats/go-multiaddr"
	"fmt"
)

func GetPeerInfo(h host.Host) peer.AddrInfo {
	return peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
}

func setupRelay() host.Host {
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/1337"),
	)
	if err != nil {
		panic(err)
	}

	// adding relay service
	_, err = relay.New(host)
	if err != nil {
		panic(err)
	}

	return host
}

func InitGameNode() {

	relay1 := setupRelay()

	fmt.Println(GetPeerInfo(relay1).ID)
}
