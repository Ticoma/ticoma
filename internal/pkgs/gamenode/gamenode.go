package gamenode

import (
	"context"
	"fmt"

	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
)

// GameNode consists of:
//
// - NetworkNode (libp2p client) for pubsub/ipfs communication (+ relay logic, if isRelay is set to true)
//
// - NodeCache to store other players' states
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

// Listens for incoming requests on the pubsub and forwards them to Cache
//
// Sends request to client after its verified
func (gn *GameNode) ListenForReqs(ctx context.Context, reqch chan interface{}) {

	var peerID string
	var data []byte

	for {
		// Listen for game requests on pubsub
		msg, err := gn.NetworkNode.Sub.Next(ctx)
		if err != nil {
			fmt.Errorf("[GAME NODE] - Error while reading msg from pubsub. Err: %s", err.Error())
		}

		peerID = msg.ReceivedFrom.String()
		data = msg.Data

		// Ignore requests from self
		if peerID == gn.NetworkNode.Host.GetPeerInfo().ID.String() {
			continue
		}

		req, err := gn.NodeCache.Put(peerID, data)
		if err != nil {
			debug.DebugLog("[GAME NODE] - Failed to process request. Err: "+err.Error(), debug.NETWORK)
		}

		// Once verified, send req to client
		reqch <- req
	}
}

func (gn *GameNode) SendRequest(ctx context.Context, data *[]byte) error {
	err := gn.NetworkNode.Topic.Publish(ctx, *data)
	if err != nil {
		return fmt.Errorf("[GAME NODE] - Failed to send request. err: %s", err.Error())
	}
	debug.DebugLog("[GAME NODE] - Request sent: "+string(*data), debug.NETWORK)
	return nil
}
