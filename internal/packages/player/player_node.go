package player

import (
	"context"
	"fmt"
	"ticoma/internal/packages/network/gamenode"
	"ticoma/internal/packages/player/nodecache"
)

// Player
type Player struct {
	*PlayerNode
}

func (p *Player) InitPlayer(ctx context.Context, gnc *gamenode.GameNodeConfig) {
	p.PlayerNode.InitPlayerNode(ctx, gnc)
}

func (p *Player) Move(posX int, posY int, destPosX int, destPosY int) error {
	ADPSchema := `{"playerId":0,"pubKey":"PUBKEY","pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`
	data := []any{posX, posY, destPosX, destPosY}
	pkg := fmt.Sprintf(ADPSchema, data...)
	fmt.Println("PACKAGE ", pkg)
	fmt.Println("CACHE ", p.PlayerNode.NodeCache)
	err := p.PlayerNode.NodeCache.Put([]byte(pkg))
	if err != nil {
		return err
	} else {
		fmt.Printf("[MOVE] Player move verified. Request: pos: {X: %d, Y: %d}, destPos: {X: %d, Y: %d}", data...)
		return nil
	}
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
