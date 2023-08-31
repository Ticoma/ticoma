package node

import (
	"context"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node/host"
	"ticoma/internal/pkgs/gamenode/network/utils"

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
	Port      string
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

	err := nn.Host.SetupHost("0.0.0.0", nodeConfig.Port)
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

	// Storage (ipfs filesystem) tests
	// h := nn.Host.GetHost()
	// s := storage.New()
	// s.StartServer(ctx, h)
	// s.Add("Hello")
	// s.Add("World!")
	// // fmt.Println(cid, err)
	// data, _ := s.GetLocal(ctx)
	// fmt.Println(data)

	nn.Topic = topic
	nn.Sub = sub
	nn.isRelay = isRelay
}
