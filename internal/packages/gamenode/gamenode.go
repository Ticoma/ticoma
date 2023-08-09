package gamenode

import (
	"context"
	"ticoma/internal/debug"

	"ticoma/internal/packages/gamenode/cache"
	"ticoma/internal/packages/gamenode/network/libp2p/node"
)

// GameNode consists of:
// - NetworkNode (libp2p client) for pubsub/ipfs communication
// - NodeCache to store own and other's package cache
type GameNode struct {
	*node.NetworkNode
	*cache.NodeCache
}

func New() *GameNode {
	nn := node.New()
	nc := cache.New()
	return &GameNode{
		NetworkNode: nn,
		NodeCache:   nc,
	}
}

func (gn *GameNode) Init(ctx context.Context, isRelay bool, nodeConfig *node.NodeConfig) {
	gn.NetworkNode.Init(ctx, isRelay, nodeConfig)
}

// Listens for incoming packages on the pubsub network, and verifies each message through the NodeCache verifier
func (gn *GameNode) ListenForPkgs(ctx context.Context) {
	for {
		msg, err := gn.NetworkNode.Sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		// don't echo own msgs
		if msg.ReceivedFrom != gn.NetworkNode.Host.GetPeerInfo().ID {
			debug.DebugLog("[PLAYER NODE] - Peer "+msg.ReceivedFrom.Pretty()+" : "+string(msg.Message.Data), debug.NETWORK)
			err := gn.NodeCache.Put(msg.Message.Data)
			if err != nil {
				debug.DebugLog("[PLAYER NODE] - I couldn't verify a package coming from "+msg.ReceivedFrom.Pretty()+"\nPkg: "+string(msg.Message.Data)+"\n "+err.Error(), debug.NETWORK)
			} else {
				debug.DebugLog("\n[PLAYER NODE] - Pkg from: "+msg.ReceivedFrom.Pretty()+" Verified !!! \n", debug.NETWORK)
			}
		} else {
			debug.DebugLog("[PLAYER NODE] - I just sent a package: "+string(msg.Message.Data), debug.NETWORK)
		}
	}
}

func (gn *GameNode) SendPkg(ctx context.Context, pkg []byte) {
	gn.NetworkNode.Topic.Publish(ctx, pkg)
}
