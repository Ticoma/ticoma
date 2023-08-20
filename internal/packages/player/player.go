package player

import (
	"context"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/packages/gamenode"
	"ticoma/internal/packages/gamenode/cache"
	"ticoma/internal/packages/gamenode/cache/interfaces"
	"ticoma/internal/packages/gamenode/network/libp2p/node"
)

// Player interface
//
// This interface is sent to client module once a connection to pubsub is established
// the client is able to perform certain actions using this interface and is limited to its functions
type Player interface {
	GetId() int
	GetPeerID() string
	GetCache() *cache.NodeCache
	Move(posX int, posY int, destPosX int, destPosY int) error
	Init(ctx context.Context, isRelay bool, nodeConfig *node.NodeConfig)
	GetPos() *interfaces.Position
	// GetPlayersPos() *map[int]interfaces.Position
}

type player struct {
	id int
	*gamenode.GameNode
	ctx context.Context
}

func New(ctx context.Context, id int) Player {
	gn := gamenode.New()
	return &player{
		id:       id,
		GameNode: gn,
		ctx:      ctx,
	}
}

func (p *player) Init(ctx context.Context, isRelay bool, nodeConfig *node.NodeConfig) {
	p.GameNode.Init(ctx, isRelay, nodeConfig)
	go p.GameNode.ListenForPkgs(ctx)
}

func (p *player) Move(posX int, posY int, destPosX int, destPosY int) error {
	ADPSchema := `{"playerId":%d,"pubKey":"PUBKEY","pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`
	data := []any{p.id, posX, posY, destPosX, destPosY}
	pkg := fmt.Sprintf(ADPSchema, data...)
	debug.DebugLog("[MOVE] PACKAGE "+pkg, debug.PLAYER)
	debug.DebugLog("[MOVE] CACHE "+fmt.Sprintf("%v", p.GameNode.NodeCache), debug.PLAYER)
	err := p.GameNode.NodeCache.Put([]byte(pkg))
	if err != nil {
		return err
	} else {
		msg := fmt.Sprintf("[MOVE] Player move verified. Id: %d, pos: {X: %d, Y: %d}, destPos: {X: %d, Y: %d}\n", data...)
		debug.DebugLog(msg, debug.PLAYER)
		// fmt.Println(msg) // DEBUG
		p.GameNode.SendPkg(p.ctx, []byte(pkg))
		return nil
	}
}

func (p *player) GetPos() *interfaces.Position {
	pos := p.NodeCache.GetCurrent(p.id).Position
	return pos
}

// Deprecated
// func (p *player) GetPlayersPos() *map[int]interfaces.Position {
// 	s := p.NodeCache.GetAll()
// 	pos := make(map[int]interfaces.Position)
// 	for id, data := range s {
// 		pos[id] = *data[1].Position
// 	}
// 	return &pos
// }

func (p *player) GetCache() *cache.NodeCache {
	return p.GameNode.NodeCache
}

func (p *player) GetPeerID() string {
	return p.GameNode.Host.GetPeerInfo().ID.String()
}

func (p *player) GetId() int {
	return p.id
}
