package gamenode

import (
	"context"
	"fmt"

	// "strings"
	// "ticoma/internal/debug"
	"ticoma/types"

	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
)

// GameNode consists of:
//
// - NetworkNode (libp2p client) for pubsub/ipfs communication (+ relay logic, if isRelay is set to true)
//
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
// Also forwards chat messages to their own channel
func (gn *GameNode) ListenForPkgs(ctx context.Context, cc chan types.ChatMessage) {
	for {

		// detect message on topic
		msg, err := gn.NetworkNode.Sub.Next(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(msg.Data))

		// listen for chat messages (CHAT_ prefix)
		// if strings.HasPrefix(string(msg.Data), "CHAT_") {
		// 	chatPkg, err := gn.NodeVerifier.SecurityVerifier.ConstructChatPkg(msg.Data)
		// 	if err != nil {
		// 		debug.DebugLog("[PLAYER NODE] - Coulnd't verify chat pkg. Err: "+err.Error(), debug.NETWORK)
		// 	}
		// 	cc <- chatPkg
		// }

		// listen for game packages
		// if msg.ReceivedFrom != gn.NetworkNode.Host.GetPeerInfo().ID { // don't echo own msgs
		// 	debug.DebugLog("[PLAYER NODE] - Peer "+msg.ReceivedFrom.Pretty()+" : "+string(msg.Message.Data), debug.NETWORK)
		// 	err := gn.NodeCache.Put(msg.Message.Data)
		// 	if err != nil {
		// 		debug.DebugLog("[PLAYER NODE] - I couldn't verify a package coming from "+msg.ReceivedFrom.Pretty()+"\nPkg: "+string(msg.Message.Data)+"\n "+err.Error(), debug.NETWORK)
		// 	} else {
		// 		debug.DebugLog("\n[PLAYER NODE] - Pkg from: "+msg.ReceivedFrom.Pretty()+" Verified !!! \n", debug.NETWORK)
		// 	}
		// } else {
		// 	debug.DebugLog("[PLAYER NODE] - I just sent a package: "+string(msg.Message.Data), debug.NETWORK)
		// }
	}
}

// Publishes ADP package to topic
func (gn *GameNode) SendADPPkg(ctx context.Context, pkg []byte) error {
	err := gn.NodeCache.Put(pkg)
	if err != nil {
		return fmt.Errorf("[PLAYER NODE] - I couldn't verify my own package!: " + err.Error())
	}
	gn.NetworkNode.Topic.Publish(ctx, pkg)
	return nil
}

// Publish chat msg to topic
func (gn *GameNode) SendChatMsg(ctx context.Context, pkg []byte) error {
	_, err := gn.SecurityVerifier.ConstructChatPkg(pkg)
	if err != nil {
		return fmt.Errorf("[PLAYER NODE] - Couldn't verify chat pkg, err: " + err.Error())
	}

	gn.NetworkNode.Topic.Publish(ctx, pkg)
	return nil
}
