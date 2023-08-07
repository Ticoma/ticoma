package core

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
)

//
// This is the core logical component of a GameNode,
// Which contains a Libp2p node (a host)
// And is used as a connector between the network and our engine,
// It is responsible for sending and receiving data on the GossipSub
//

type GameNodeCore struct {
	host                  host.Host // Basic libp2p client
	relayConnectionStatus bool      // Is node core connected to a Relay (Standalone GameNode)
}

// Initialize host
func (gnc *GameNodeCore) SetupHost(listenIp string, listenPort string) error {

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + listenIp + "/tcp/" + listenPort),
	)

	if err != nil {
		return err
	}

	gnc.host = h

	return nil
}

// Peer info
func (gnc *GameNodeCore) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    gnc.host.ID(),
		Addrs: gnc.host.Addrs(),
	}
}

// Send reserve slot request to relay
func (gnc *GameNodeCore) ReserveSlot(ctx context.Context, relayAddrInfo peer.AddrInfo) {
	_, err := client.Reserve(ctx, gnc.host, relayAddrInfo)
	if err != nil {
		panic(err)
	}
}

// Establish connection to relay
func (gnc *GameNodeCore) ConnectToRelay(ctx context.Context, relayAddrInfo peer.AddrInfo) {

	err := gnc.host.Connect(ctx, relayAddrInfo)
	if err != nil {
		panic(err)
	}

	// Connection to relay established
	gnc.relayConnectionStatus = true
}

func (gnc *GameNodeCore) TEST(ctx *context.Context) {
	fmt.Println(ctx, "I GOT CONTEXT")
}

// Connect to pubsub and return topic, sub
func (gnc *GameNodeCore) ConnectToPubsub(ctx context.Context, topicName string, relayEnabled bool) (*pubsub.Topic, *pubsub.Subscription) {

	ps, err := pubsub.NewGossipSub(ctx, gnc.host)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join(topicName)
	if err != nil {
		panic(err)
	}

	if relayEnabled {
		topic.Relay()
	}

	sub, subErr := topic.Subscribe()
	if subErr != nil {
		panic(err)
	}

	fmt.Printf("Connected to topic: %s!\n", topicName)
	return topic, sub
}