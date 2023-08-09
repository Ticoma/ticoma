package player

import (
	"context"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/packages/gamenode"
	"ticoma/internal/packages/gamenode/network/libp2p/node"
)

// Player interface
//
// This interface is sent to client module once a connection to pubsub is established
// the client is able to perform certain actions using this interface and is limited to its functions
type Player interface {
	GetPeerID() string
	Move(posX int, posY int, destPosX int, destPosY int) error
	Init(ctx context.Context, isRelay bool, nodeConfig *node.NodeConfig)
}

type player struct {
	*gamenode.GameNode
	ctx context.Context
}

func New(ctx context.Context) Player {
	gn := gamenode.New()
	return &player{
		GameNode: gn,
		ctx:      ctx,
	}
}

func (p *player) Init(ctx context.Context, isRelay bool, nodeConfig *node.NodeConfig) {
	p.GameNode.Init(ctx, isRelay, nodeConfig)
}

func (p *player) Move(posX int, posY int, destPosX int, destPosY int) error {
	ADPSchema := `{"playerId":0,"pubKey":"PUBKEY","pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`
	data := []any{posX, posY, destPosX, destPosY}
	pkg := fmt.Sprintf(ADPSchema, data...)
	debug.DebugLog("[MOVE] PACKAGE "+pkg, debug.PLAYER)
	debug.DebugLog("[MOVE] CACHE "+fmt.Sprintf("%v", p.GameNode.NodeCache), debug.PLAYER)
	err := p.GameNode.NodeCache.Put([]byte(pkg))
	if err != nil {
		return err
	} else {
		msg := fmt.Sprintf("[MOVE] Player move verified. Request: pos: {X: %d, Y: %d}, destPos: {X: %d, Y: %d}", data...)
		debug.DebugLog(msg, debug.PLAYER)
		fmt.Println(msg)
		p.GameNode.SendPkg(p.ctx, []byte(pkg))
		return nil
	}
}

func (p *player) GetPeerID() string {
	return p.GameNode.Host.GetPeerInfo().ID.String()
}
