package gamenode

import (
	"context"
	"ticoma/internal/debug"
	"ticoma/internal/packages/network/gamenode/core"
	"ticoma/internal/packages/network/utils"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// GameNode, integral part of PlayerNode
type GameNode struct {
	*core.GameNodeCore
	Topic *pubsub.Topic
	Sub   *pubsub.Subscription
}

// Config fields needed to initialize a functional GameNode
type NodeConfig struct {
	RelayAddr string // In GameNode -> addr/ip/port of relay to connect through, in Standalone -> init data
	RelayIp   string
	RelayPort string
	IsRelay   bool
}

func New() *GameNode {
	return &GameNode{
		GameNodeCore: &core.GameNodeCore{},
	}
}

// func NewStandaloneGameNode(ctx context.Context, nodeConfig *NodeConfig) *StandaloneGameNode {
// 	ign := NewIntegralGameNode()
// 	ign.InitIntegralGameNode(ctx, nodeConfig)
// 	return &StandaloneGameNode{
// 		IntegralGameNode: ign,
// 		GameNodeRelay:    &relay.GameNodeRelay{},
// 	}
// }

// Connects the Integral Game Node to PubSub
func (gn *GameNode) InitGameNode(ctx context.Context, nodeConfig *NodeConfig) {

	err := gn.GameNodeCore.SetupHost("0.0.0.0", "1337")
	if err != nil {
		panic(err)
	}

	if !nodeConfig.IsRelay {
		relayInfo := utils.ConvertToAddrInfo(nodeConfig.RelayIp, nodeConfig.RelayAddr, nodeConfig.RelayPort)
		gn.ConnectToRelay(ctx, *relayInfo)
		debug.DebugLog("Connected to relay!", debug.NETWORK)
	}

	topic, sub := gn.ConnectToPubsub(ctx, "ticoma1", nodeConfig.IsRelay)
	debug.DebugLog("Connected to pubsub!", debug.NETWORK)

	gn.Topic = topic
	gn.Sub = sub
}
