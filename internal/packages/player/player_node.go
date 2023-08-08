package player

import (
	"context"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/packages/network/gamenode"
	"ticoma/internal/packages/player/nodecache"
)

// Player interface
//
// This interface is sent to client module once a connection to pubsub is established
// the client is able to perform certain actions through this interface
type PlayerInterface interface {
	GetPeerID() string
	Move(posX int, posY int, destPosX int, destPosY int) error
}

type Player struct {
	*PlayerNode
	ctx context.Context
}

func New(ctx context.Context, nodeConfig *gamenode.NodeConfig) PlayerInterface {
	pn := NewPlayerNode()
	pn.InitPlayerNode(ctx, nodeConfig)
	return &Player{
		PlayerNode: pn,
		ctx:        ctx,
	}
}

func (p *Player) Move(posX int, posY int, destPosX int, destPosY int) error {
	ADPSchema := `{"playerId":0,"pubKey":"PUBKEY","pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`
	data := []any{posX, posY, destPosX, destPosY}
	pkg := fmt.Sprintf(ADPSchema, data...)
	debug.DebugLog("[MOVE] PACKAGE "+pkg, debug.PLAYER)
	debug.DebugLog("[MOVE] CACHE "+fmt.Sprintf("%v", p.PlayerNode.NodeCache), debug.PLAYER)
	err := p.PlayerNode.NodeCache.Put([]byte(pkg))
	if err != nil {
		return err
	} else {
		msg := fmt.Sprintf("[MOVE] Player move verified. Request: pos: {X: %d, Y: %d}, destPos: {X: %d, Y: %d}", data...)
		debug.DebugLog(msg, debug.PLAYER)
		fmt.Println(msg)
		p.PlayerNode.SendPkg(p.ctx, pkg)
		return nil
	}
}

func (p *Player) GetPeerID() string {
	return p.PlayerNode.GetPeerInfo().ID.String()
}

// PlayerNode consists of:
// - GameNode mechanism which allows direct communication with libp2p/ipfs network, pubsub etc.
// - NodeCache to store own and other's package cache
// - Middleware for quick & handy signal handling from the Client

type PlayerNode struct {
	*gamenode.GameNode
	*nodecache.NodeCache
}

func NewPlayerNode() *PlayerNode {
	nc := nodecache.New()
	gn := gamenode.New()
	return &PlayerNode{
		GameNode:  gn,
		NodeCache: nc,
	}
}

func (pn *PlayerNode) InitPlayerNode(ctx context.Context, nodeConfig *gamenode.NodeConfig) {
	pn.GameNode.InitGameNode(ctx, nodeConfig)
	go pn.ListenForPkgs(ctx)
}

// Listens for incoming packages on the pubsub network, and verifies each message through the NodeCache verifier
func (pn *PlayerNode) ListenForPkgs(ctx context.Context) {
	for {
		msg, err := pn.GameNode.Sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		// don't echo own msgs
		if msg.ReceivedFrom != pn.GameNode.GetPeerInfo().ID {
			debug.DebugLog("[PLAYER NODE] - Peer "+msg.ReceivedFrom.Pretty()+" : "+string(msg.Message.Data), debug.NETWORK)
			err := pn.NodeCache.Put(msg.Message.Data)
			if err != nil {
				debug.DebugLog("[PLAYER NODE] - I couldn't verify a package coming from "+msg.ReceivedFrom.Pretty()+"\nPkg: "+string(msg.Message.Data), debug.NETWORK)
			}
		} else {
			debug.DebugLog("[PLAYER NODE] - I just sent a package: "+string(msg.Message.Data), debug.NETWORK)
		}
	}
}

func (pn *PlayerNode) SendPkg(ctx context.Context, data string) {
	pn.GameNode.Topic.Publish(ctx, []byte(data))
}
