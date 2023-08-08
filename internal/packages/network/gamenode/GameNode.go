package gamenode

import (
	"context"
	"ticoma/internal/debug"
	"ticoma/internal/packages/network/gamenode/core"
	"ticoma/internal/packages/network/utils"

	relay "ticoma/internal/packages/network/gamenode/relay"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Extended game node,
// Consists of: GameNodeCore + NodeRelay
type StandaloneGameNode struct {
	*IntegralGameNode
	*relay.GameNodeRelay
}

// GameNode, integral part of PlayerNode
type IntegralGameNode struct {
	*core.GameNodeCore
	Topic *pubsub.Topic
	Sub   *pubsub.Subscription
}

// Config fields needed to initialize a functional GameNode
type GameNodeConfig struct {
	RelayAddr          string // In GameNode -> addr/ip/port of relay to connect through, in Standalone -> init data
	RelayIp            string
	RelayPort          string
	EnableDebugLogging bool
}

func NewIntegralGameNode() *IntegralGameNode {
	return &IntegralGameNode{
		GameNodeCore: &core.GameNodeCore{},
	}
}

func NewStandaloneGameNode() *StandaloneGameNode {
	ign := NewIntegralGameNode()
	return &StandaloneGameNode{
		IntegralGameNode: ign,
		GameNodeRelay:    &relay.GameNodeRelay{},
	}
}

// Initializes an empty GameNode in a configurable way
func (ign *IntegralGameNode) InitGameNode(ctx context.Context, config *GameNodeConfig) {

	relayInfo := utils.ConvertToAddrInfo(config.RelayIp, config.RelayAddr, config.RelayPort)

	// Host setup
	err := ign.GameNodeCore.SetupHost("0.0.0.0", "8888")
	if err != nil {
		panic(err)
	}

	debug.DebugLog("GameNode host set up", debug.NETWORK)

	ign.ConnectToRelay(ctx, *relayInfo)
	debug.DebugLog("GameNode host set up", debug.NETWORK)

	// Pubsub
	topic, sub := ign.ConnectToPubsub(ctx, "ticoma1", true)
	debug.DebugLog("Connected to pubsub!", debug.NETWORK)

	ign.Topic = topic
	ign.Sub = sub
}
