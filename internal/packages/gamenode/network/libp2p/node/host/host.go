package host

import (
	"context"

	debug "ticoma/internal/debug"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

//
// According to libp2p lib docs:
// Host represents a single libp2p node in a peer-to-peer network.
// (https://pkg.go.dev/github.com/libp2p/go-libp2p-core/host)
//

// Libp2p peer
type Libp2pHost struct {
	host host.Host
}

func (h *Libp2pHost) SetupHost(listenIp string, listenPort string) error {

	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + listenIp + "/tcp/" + listenPort),
	)

	if err != nil {
		return err
	} else {
		h.host = host
		return nil
	}
}

func (h *Libp2pHost) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    h.host.ID(),
		Addrs: h.host.Addrs(),
	}
}

func (h *Libp2pHost) ConnectToRelay(ctx context.Context, relayAddrInfo peer.AddrInfo) {

	err := h.host.Connect(ctx, relayAddrInfo)
	if err != nil {
		panic(err)
	}
}

// Returns [*topic, *sub] on successful connection
func (h *Libp2pHost) ConnectToPubsub(ctx context.Context, topicName string, isRelay bool) (*pubsub.Topic, *pubsub.Subscription) {

	ps, err := pubsub.NewGossipSub(ctx, h.host)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join(topicName)
	if err != nil {
		panic(err)
	}

	if isRelay {
		topic.Relay()
	}

	sub, subErr := topic.Subscribe()
	if subErr != nil {
		panic(err)
	}

	debug.DebugLog("Connected to topic: "+topicName, debug.NETWORK)
	return topic, sub
}

func (h *Libp2pHost) GetHost() host.Host {
	return h.host
}
