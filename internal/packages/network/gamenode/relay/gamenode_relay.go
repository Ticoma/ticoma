package relay

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

// Public relay node service for standalone game nodes,
// Players can use those relays to join the game
type GameNodeRelay struct {
	RelayHost host.Host
}

func (gnr *GameNodeRelay) SetupRelay(listenIp string, listenPort string) {
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + listenIp + "/tcp/" + listenPort),
	)

	if err != nil {
		panic(err)
	}

	// adding relay service
	_, err = relay.New(host)
	if err != nil {
		panic(err)
	}

	gnr.RelayHost = host
}
