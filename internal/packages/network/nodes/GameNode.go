package nodes

import (
	"context"
	"fmt"
	"ticoma/packages/network/nodes/core"
	"ticoma/packages/network/utils"
)

// Basic GameNode, integral part of PlayerNode
type GameNode struct {
	*core.GameNodeCore
}

func NewGameNode() *GameNode {
	return &GameNode{
		GameNodeCore: &core.GameNodeCore{},
	}
}

// Initializes an empty GameNode in a configurable way
func (gn *GameNode) InitGameNode(ctx *context.Context, relayAddr string, relayIp string, relayPort string, enableDebugLogging bool) {

	relayInfo := utils.ConvertToAddrInfo(relayIp, relayAddr, relayPort)

	// Host setup
	gn.GameNodeCore.SetupHost("127.0.0.1", "1337")
	if enableDebugLogging {
		fmt.Println("GameNode host set up")
	}

	// Connect to relay (TODO: check if reservation is needed)
	gn.ConnectToRelay(*ctx, *relayInfo)
	if enableDebugLogging {
		fmt.Println("GameNode connected to relay")
	}

	// gn.ReserveSlot(ctx, *relayInfo)
	// fmt.Println("GameNode relay slot reserved")

	// Pubsub
	gn.ConnectToPubsub(*ctx, "ticoma1", true) //
	if enableDebugLogging {
		fmt.Println("Connected to pubsub!")
	}
}
