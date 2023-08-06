package player

import (
	"context"
	"fmt"
	"ticoma/packages/network/gamenode"
	"ticoma/packages/network/gamenode/core"
	nodecache "ticoma/packages/nodes/modules"
	"ticoma/packages/nodes/modules/verifier"
)

//
// PlayerNode
// The internal interface of a communicator between the Network and the Client
//
// Consists of:
// - GameNode mechanism which allows direct communication with libp2p/ipfs network, pubsub etc.
// - NodeCache to store own and other's package cache
// - Middleware for quick & handy signal handling from the Client

type PlayerNode struct {
	*gamenode.GameNode
	*nodecache.NodeCache
	// Middleware soon (after basic @raylib implementation)
}

func NewPlayerNode() *PlayerNode {
	return &PlayerNode{
		GameNode: &gamenode.GameNode{
			GameNodeCore: &core.GameNodeCore{},
		},
		NodeCache: &nodecache.NodeCache{
			Verifier: &verifier.Verifier{},
		},
	}
}

func (pn *PlayerNode) InitPlayerNode(ctx context.Context, gameNodeConfig *gamenode.GameNodeConfig) {
	pn.GameNode.InitGameNode(ctx, gameNodeConfig)
}

// Listens for incoming packages on the pubsub network, and verifies each message through the NodeCache verifier
func (pn *PlayerNode) ListenForPkgs(ctx context.Context) {
	for {
		msg, err := pn.GameNode.Sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.ReceivedFrom, ": ", string(msg.Message.Data))
		// pn.NodeCache.Put(msg.Message.Data)
		// TODO: Find a way to convert []byte / string -> ADPT quickly
	}
}

func (pn *PlayerNode) SendPkg(ctx context.Context, data string) {
	pn.GameNode.Topic.Publish(ctx, []byte(data))
}
