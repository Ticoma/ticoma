package core

import (
	"context"

	debug "ticoma/internal/debug"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Logical interface responsible for direct communication on the pubsub
type GameNodeCore struct {
	host host.Host
}

func (gnc *GameNodeCore) SetupHost(listenIp string, listenPort string) error {

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + listenIp + "/tcp/" + listenPort),
	)

	if err != nil {
		return err
	} else {
		gnc.host = h
		return nil
	}
}

func (gnc *GameNodeCore) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    gnc.host.ID(),
		Addrs: gnc.host.Addrs(),
	}
}

func (gnc *GameNodeCore) ConnectToRelay(ctx context.Context, relayAddrInfo peer.AddrInfo) {

	err := gnc.host.Connect(ctx, relayAddrInfo)
	if err != nil {
		panic(err)
	}
}

// returns pubsub[topic, sub] objects
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

	debug.DebugLog("Connected to topic: "+topicName, debug.NETWORK)
	return topic, sub
}
