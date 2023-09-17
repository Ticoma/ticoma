package player

import (
	"context"
	"encoding/json"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
	"ticoma/types"
)

// Player interface
//
// This interface is sent to client module once a connection to pubsub is established
// the client is able to perform certain actions using this interface and is limited to its functions
type Player interface {
	Register() error
	Login() error
	Move(posX *int, posY *int, destPosX *int, destPosY *int) error
	Chat(msg *[]byte) error
	Init(ctx context.Context, reqch chan interface{}, isRelay bool, nodeConfig *node.NodeConfig)
	GetPeerID() string
	GetCache() *cache.Memory
	GetPos() *types.PlayerPosition
}

type player struct {
	*gamenode.GameNode
	ctx context.Context
}

func New(ctx context.Context) *player {
	gn := gamenode.New()
	return &player{
		GameNode: gn,
		ctx:      ctx,
	}
}

func (p *player) Init(ctx context.Context, reqch chan interface{}, isRelay bool, nodeConfig *node.NodeConfig) {
	p.GameNode.Init(ctx, isRelay, nodeConfig)
	go p.GameNode.ListenForReqs(ctx, reqch)
	go p.GameNode.PerformSnapshot(ctx)
}

//
// Account requests
//

func (p *player) Register() error {
	// TODO: add nickname support
	pfx := []byte(security.REGISTER_PREFIX)
	err := p.SendRequest(p.ctx, &pfx)
	return err
}

func (p *player) Login() error {
	pfx := []byte(security.LOGIN_PREFIX)
	err := p.SendRequest(p.ctx, &pfx)
	return err
}

//
// Game requests
//

func (p *player) Move(posX *int, posY *int, destPosX *int, destPosY *int) error {

	pfx := []byte(security.MOVE_PREFIX)
	pos := types.Position{X: *posX, Y: *posY}
	destPos := types.DestPosition{X: *destPosX, Y: *destPosY}
	pp := &types.PlayerPosition{Position: pos, DestPosition: destPos}

	moveReqJSON, err := json.Marshal(pp)
	if err != nil {
		return fmt.Errorf("[PLAYER] - Failed to serialize request. Err: %s", err.Error())
	}
	var moveReq []byte = append(pfx, moveReqJSON...)
	p.SendRequest(p.ctx, &moveReq)

	debug.DebugLog("[MOVE] Sending move req: "+string(moveReq), debug.PLAYER)
	return nil
}

func (p *player) Chat(msg *[]byte) error {
	pfx := []byte("CHAT_")
	var reqData []byte = append(pfx, *msg...)
	err := p.SendRequest(p.ctx, &reqData)
	if err != nil {
		return fmt.Errorf("[PLAYER] - Failed to send chat message. Err: %s", err.Error())
	}
	return nil
}

//
// Getters
//

func (p *player) GetCache() *cache.Memory {
	return p.GameNode.GetAll()
}

func (p *player) GetPeerID() string {
	return p.GameNode.NetworkNode.Host.GetPeerInfo().ID.String()
}

func (p *player) GetPos() *types.PlayerPosition {
	return p.GameNode.GetCurrPlayerPos(p.GetPeerID())
}
