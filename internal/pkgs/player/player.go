package player

import (
	"context"
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
	Register(nickname *string) error
	Login() error
	Move(posX *int, posY *int, destPosX *int, destPosY *int) error
	Chat(msg *[]byte) error
	Init(ctx context.Context, crc chan types.CachedRequest, isRelay bool, nodeConfig *node.NodeConfig)
	GetPeerID() *string
	GetNickname(peerID *string) *string
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

func (p *player) Init(ctx context.Context, crc chan types.CachedRequest, isRelay bool, nodeConfig *node.NodeConfig) {
	p.GameNode.Init(ctx, isRelay, nodeConfig)
	go p.GameNode.ListenForReqs(ctx, crc)
	go p.GameNode.PerformSnapshot(ctx)
}

//
// Account requests
//

func (p *player) Register(nickname *string) error {
	pfx := []byte(security.REGISTER_PREFIX)
	reqData := append(pfx, []byte(*nickname)...)
	err := p.SendRequest(p.ctx, &reqData)
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
	moveReq := []byte(fmt.Sprintf(`MOVE_{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, *posX, *posY, *destPosX, *destPosY))
	err := p.SendRequest(p.ctx, &moveReq)
	if err != nil {
		return fmt.Errorf("[PLAYER] - Failed to send move request. Err: %s", err.Error())
	}
	debug.DebugLog("[MOVE] Sending move req: "+string(moveReq), debug.PLAYER)
	return nil
}

func (p *player) Chat(msg *[]byte) error {
	pfx := []byte("CHAT_")
	reqData := fmt.Sprintf(`{"message":"%s"}`, string(*msg))
	chatReq := append(pfx, []byte(reqData)...)
	err := p.SendRequest(p.ctx, &chatReq)
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

func (p *player) GetPeerID() *string {
	peerID := p.GameNode.NetworkNode.Host.GetPeerInfo().ID.String()
	return &peerID
}

func (p *player) GetPos() *types.PlayerPosition {
	return p.GameNode.GetCurrPlayerPos(*p.GetPeerID())
}

func (p *player) GetNickname(peerID *string) *string {
	return p.GameNode.GetNickname(*peerID)
}
