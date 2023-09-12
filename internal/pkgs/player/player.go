package player

import (
	"context"
	"encoding/json"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
	"ticoma/types"
)

// Player interface
//
// This interface is sent to client module once a connection to pubsub is established
// the client is able to perform certain actions using this interface and is limited to its functions
type Player interface {
	GetPeerID() string
	GetCache() *cache.NodeCache
	Move(posX *int, posY *int, destPosX *int, destPosY *int) error
	Chat(msg *[]byte)
	Init(ctx context.Context, reqch chan interface{}, isRelay bool, nodeConfig *node.NodeConfig)
	GetPos() *types.Position
}

type player struct {
	*gamenode.GameNode
	ctx context.Context
}

func New(ctx context.Context, id int) *player {
	gn := gamenode.New()
	return &player{
		GameNode: gn,
		ctx:      ctx,
	}
}

func (p *player) Init(ctx context.Context, reqch chan interface{}, isRelay bool, nodeConfig *node.NodeConfig) {
	p.GameNode.Init(ctx, isRelay, nodeConfig)
	go p.GameNode.ListenForReqs(ctx, reqch)
}

//
// Request-related funcs
//

func (p *player) Move(posX *int, posY *int, destPosX *int, destPosY *int) error {

	prefix := []byte("MOVE_")
	pos := types.Position{X: *posX, Y: *posY}
	destPos := types.DestPosition{X: *destPosX, Y: *destPosY}
	pp := &types.PlayerPosition{Position: &pos, DestPosition: &destPos}

	moveReqJSON, err := json.Marshal(pp)
	if err != nil {
		return fmt.Errorf("[PLAYER] - Failed to serialize request. Err: %v", err)
	}
	var moveReq []byte = append(prefix, moveReqJSON...)
	fmt.Println("MOVE REQ: ", string(moveReq)) //tmp
	p.SendRequest(p.ctx, &moveReq)

	debug.DebugLog("[MOVE] Sending move req: "+string(moveReq), debug.PLAYER)
	return nil
}

func (p *player) Chat(msg *[]byte) {
	prefix := []byte("CHAT_")
	var chatReq []byte = append(prefix, *msg...)
	p.SendRequest(p.ctx, &chatReq)
}

//
// Getters
//

func (p *player) GetCache() *cache.NodeCache {
	return p.GetCache()
}

func (p *player) GetPeerID() string {
	return p.GetPeerID()
}

func (p *player) GetPos() *types.Position {
	return p.GetCurrPlayerPos(p.GameNode.Host.GetPeerInfo().String())
}
