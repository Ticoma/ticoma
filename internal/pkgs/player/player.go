package player

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetId() int
	GetPeerID() string
	GetCache() *cache.NodeCache
	Move(posX int, posY int, destPosX int, destPosY int) error
	Chat(msg []byte) error
	Init(ctx context.Context, cc chan types.ChatMessage, isRelay bool, nodeConfig *node.NodeConfig)
	// GetPos() *interfaces.Position
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

func (p *player) Move(posX int, posY int, destPosX int, destPosY int) error {

	// "MOVE" Prefix
	prefix := "MOVE_"
	pos := types.Position{
		X: posX,
		Y: posY,
	}
	destPos := types.DestPosition{
		X: destPosX,
		Y: destPosY,
	}
	pp := &types.PlayerPosition{
		Position:     &pos,
		DestPosition: &destPos,
	}

	moveReqJSON, err := json.Marshal(pp)
	if err != nil {
		return fmt.Errorf("[PLAYER] - Failed to serialize request. Err: ", err)
	}
	moveReq := prefix + string(moveReqJSON)
	fmt.Println(moveReq)
	return nil

	// debug.DebugLog("[MOVE] PACKAGE "+pkg, debug.PLAYER)
	// debug.DebugLog("[MOVE] CACHE "+fmt.Sprintf("%v", p.GameNode.NodeCache), debug.PLAYER)
}
