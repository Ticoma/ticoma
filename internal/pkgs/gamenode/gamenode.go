package gamenode

import (
	"context"
	"fmt"

	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
	"ticoma/types"
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

// Listen for incoming requests on Topic and put them in Cache.
// Forwards request to client if it gets verified
func (gn *GameNode) ListenForReqs(ctx context.Context, crc chan types.CachedRequest) {

	var peerID string
	var data []byte

	for {
		// Listen for game requests on pubsub
		msg, err := gn.NetworkNode.Sub.Next(ctx)
		if err != nil {
			debug.DebugLog(fmt.Sprintf("[GAME NODE] - Error while reading msg from pubsub. Err: %s", err.Error()), debug.NETWORK)
		}

		peerID = msg.ReceivedFrom.String()
		data = msg.Data

		req, pfx, err := gn.NodeCache.Put(peerID, data)
		if err != nil {
			debug.DebugLog("[GAME NODE] - Failed to process request. Err: "+err.Error(), debug.NETWORK)
			continue
		}

		// If valid req, forward to client
		if req != nil && pfx != "" {
			cachedReq := types.CachedRequest{
				Pfx: pfx,
				Req: req,
			}
			crc <- cachedReq
		}
	}
}

// Send a request to the network and Put in own cache
func (gn *GameNode) SendRequest(ctx context.Context, data *[]byte) error {
	err := gn.Topic.Publish(ctx, *data)
	if err != nil {
		return fmt.Errorf("[GAME NODE] - Failed to send request. Err: %s", err.Error())
	}
	return nil
}
