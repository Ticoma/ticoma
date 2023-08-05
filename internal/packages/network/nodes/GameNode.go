package nodes

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	 "context"
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


func connectToPubsub(ctx context.Context, h host.Host) (*pubsub.Topic, *pubsub.Subscription) {

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join("ticoma2")
	if err != nil {
		panic(err)
	}
	topic.Relay()

	sub, subErr := topic.Subscribe()
	if subErr != nil {
		panic(err)
	}

	fmt.Println("Connected to pubsub!")
	return topic, sub
}

func InitGameNode() {
	ctx := context.Background()
	relay1 := setupRelay()	
	connectToPubsub(ctx, relay1)
	fmt.Println(GetPeerInfo(relay1).ID)
}
