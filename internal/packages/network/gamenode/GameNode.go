package gamenode

import (
	"context"
	"fmt"
	"ticoma/internal/packages/network/gamenode/core"
	"ticoma/internal/packages/network/utils"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

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

// Initializes an empty GameNode in a configurable way
func (ign *IntegralGameNode) InitGameNode(ctx context.Context, config *GameNodeConfig) {

	relayInfo := utils.ConvertToAddrInfo(config.RelayIp, config.RelayAddr, config.RelayPort)

	// Host setup
	err := ign.GameNodeCore.SetupHost("0.0.0.0", "8888")
	if err != nil {
		panic(err)
	}

	if config.EnableDebugLogging {
		fmt.Println("GameNode host set up")
	}

	ign.ConnectToRelay(ctx, *relayInfo)
	if config.EnableDebugLogging {
		fmt.Println("GameNode connected to relay")
	}

	// Pubsub
	topic, sub := ign.ConnectToPubsub(ctx, "ticoma1", true)
	if config.EnableDebugLogging {
		fmt.Println("Connected to pubsub!")
	}

	ign.Topic = topic
	ign.Sub = sub
}
