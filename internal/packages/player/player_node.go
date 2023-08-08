package player

import (
	"context"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/packages/network/gamenode"
	"ticoma/internal/packages/player/nodecache"
)

// Player interface
type PlayerInterface interface {
	GetPeerID() string
	Move(posX int, posY int, destPosX int, destPosY int) error
}

// Player
type Player struct {
	*PlayerNode
}

func New(ctx context.Context, gnc *gamenode.GameNodeConfig) PlayerInterface {
	pn := NewPlayerNode()
	pn.InitPlayerNode(ctx, gnc)
	return &Player{
		PlayerNode: pn,
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
	*gamenode.IntegralGameNode
	*nodecache.NodeCache
}

func NewPlayerNode() *PlayerNode {
	nc := nodecache.New()
	ign := gamenode.NewIntegralGameNode()
	return &PlayerNode{
		IntegralGameNode: ign,
		NodeCache:        nc,
	}
}

func (pn *PlayerNode) InitPlayerNode(ctx context.Context, gameNodeConfig *gamenode.GameNodeConfig) {
	pn.IntegralGameNode.InitGameNode(ctx, gameNodeConfig)
	go pn.ListenForPkgs(ctx)
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
		} else {
			debug.DebugLog("[PLAYER NODE] - I just sent a package: "+string(msg.Message.Data), debug.NETWORK)
		}
	}
}

func (pn *PlayerNode) SendPkg(ctx context.Context, data string) {
	pn.IntegralGameNode.Topic.Publish(ctx, []byte(data))
}
