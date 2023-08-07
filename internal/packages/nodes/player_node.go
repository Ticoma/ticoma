package player

import (
	"context"
	"fmt"
	"ticoma/internal/packages/network/gamenode"
	"ticoma/internal/packages/network/gamenode/core"
	nodecache "ticoma/internal/packages/nodes/modules"
	"ticoma/internal/packages/nodes/modules/verifier"
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
	*gamenode.IntegralGameNode
	*nodecache.NodeCache
}

func NewPlayerNode() *PlayerNode {
	return &PlayerNode{
		IntegralGameNode: &gamenode.IntegralGameNode{
			GameNodeCore: &core.GameNodeCore{},
		},
		NodeCache: &nodecache.NodeCache{
			NodeVerifier: &verifier.NodeVerifier{},
		},
	}
}

func (pn *PlayerNode) InitPlayerNode(ctx context.Context, gameNodeConfig *gamenode.GameNodeConfig) {
	pn.IntegralGameNode.InitGameNode(ctx, gameNodeConfig)
}

// Listens for incoming packages on the pubsub network, and verifies each message through the NodeCache verifier
func (pn *PlayerNode) ListenForPkgs(ctx context.Context) {
	for {
		msg, err := pn.IntegralGameNode.Sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		// don't echo own msgs
		if msg.ReceivedFrom != pn.GetPeerInfo().ID {
			fmt.Println(msg.ReceivedFrom, ": ", string(msg.Message.Data))
			// pn.NodeCache.Put(msg.Message.Data) <-- soon
		}
	}
}

func (pn *PlayerNode) SendPkg(ctx context.Context, data string) {
	pn.IntegralGameNode.Topic.Publish(ctx, []byte(data))
}
