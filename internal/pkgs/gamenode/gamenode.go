package gamenode

import (
	"context"

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
			panic(err)
		}

		peerID = msg.ReceivedFrom.String()
		data = msg.Data

		// Ignore requests from self
		if peerID == gn.NetworkNode.Host.GetPeerInfo().ID.String() {
			continue
		}

		req, err := gn.NodeCache.Put(peerID, data)
		if err != nil {
			debug.DebugLog("[PLAYER NODE] - Failed to process request. Err: "+err.Error(), debug.NETWORK)
		}

		// Once verified, send req to client
		reqch <- req
	}
}

func (gn *GameNode) SendRequest(ctx context.Context, data *[]byte) {
	gn.NetworkNode.Topic.Publish(ctx, *data)
	debug.DebugLog("[PLAYER NODE] - I just sent a request: "+string(*data), debug.NETWORK)
}
