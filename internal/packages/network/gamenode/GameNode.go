package gamenode

import (
	"context"
	"fmt"
	"ticoma/packages/network/gamenode/core"
	"ticoma/packages/network/utils"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Basic GameNode, integral part of PlayerNode
type GameNode struct {
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

func NewGameNode() *GameNode {
	return &GameNode{
		GameNodeCore: &core.GameNodeCore{},
	}
}

// Initializes an empty GameNode in a configurable way
func (gn *GameNode) InitGameNode(ctx context.Context, config *GameNodeConfig) {

	relayInfo := utils.ConvertToAddrInfo(config.RelayIp, config.RelayAddr, config.RelayPort)

	// Host setup
	err := gn.GameNodeCore.SetupHost("0.0.0.0", "8888")
	if err != nil {
		panic(err)
	}

	if config.EnableDebugLogging {
		fmt.Println("GameNode host set up")
	}

	// Connect to relay (TODO: check if reservation is needed)
	// gn.ConnectToRelay(config.Ctx, *relayInfo)
	gn.ConnectToRelay(ctx, *relayInfo)
	if config.EnableDebugLogging {
		fmt.Println("GameNode connected to relay")
	}

	// gn.ReserveSlot(ctx, *relayInfo)
	// fmt.Println("GameNode relay slot reserved")

	// Pubsub
	topic, sub := gn.ConnectToPubsub(ctx, "ticoma1", true) //
	// topic, sub := gn.ConnectToPubsub(config.Ctx, "ticoma1", true) //
	if config.EnableDebugLogging {
		fmt.Println("Connected to pubsub!")
	}

	gn.Topic = topic
	gn.Sub = sub
}
