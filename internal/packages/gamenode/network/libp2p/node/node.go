package node

import (
	"context"
	"ticoma/internal/debug"
	"ticoma/internal/packages/gamenode/network/libp2p/host"
	"ticoma/internal/packages/gamenode/network/utils"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type NetworkNode struct {
	Host    host.Libp2pHost
	Topic   *pubsub.Topic
	Sub     *pubsub.Subscription
	isRelay bool
}

// Config opts for NetworkNode initialization
type NodeConfig struct {
	RelayAddr string
	RelayIp   string
	RelayPort string
}

func New() *NetworkNode {
	return &NetworkNode{}
}

// Init NetworkNode and connect to pubsub
//
// To set up a relay NetworkNode (Standalone), isRelay must be set to true
func (nn *NetworkNode) Init(ctx context.Context, isRelay bool, nodeConfig *NodeConfig) {

	err := nn.Host.SetupHost("0.0.0.0", "1337")
	if err != nil {
		panic(err)
	}

	if !isRelay {
		relayInfo := utils.ConvertToAddrInfo(nodeConfig.RelayIp, nodeConfig.RelayAddr, nodeConfig.RelayPort)
		nn.Host.ConnectToRelay(ctx, *relayInfo)
		debug.DebugLog("Connected to relay!", debug.NETWORK)
	}

	topic, sub := nn.Host.ConnectToPubsub(ctx, "ticoma1", isRelay)
	debug.DebugLog("Connected to pubsub!", debug.NETWORK)

	nn.Topic = topic
	nn.Sub = sub
	nn.isRelay = isRelay
}
