package cache

import (
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache/verifier"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
)

type NodeCache struct {
	Memory
	*verifier.NodeVerifier
}

type PlayerStates struct {
	Prev *types.Player
	Curr *types.Player
}

type Memory map[string]PlayerStates

func New() *NodeCache {
	v := verifier.New()
	return &NodeCache{
		Memory:       Memory{},
		NodeVerifier: v,
	}
}

func (nc *NodeCache) GetAll() Memory {
	return nc.Memory
}

func (nc *NodeCache) GetPlayer(peerID string) PlayerStates {
	return nc.Memory[peerID]
}

func (nc *NodeCache) GetPrevPlayerPos(peerID string) *types.Player {
	return nc.Memory[peerID].Prev
}

func (nc *NodeCache) GetCurrPlayerPos(peerID string) *types.Player {
	return nc.Memory[peerID].Curr
}

// Put new data to NodeCache
//
// Returns automatically constructed request interface (based on data preifx)
func (nc *NodeCache) Put(peerID string, data []byte) (interface{}, error) {

	// TODO : Add instant req decline if !nc.Memory[peerID].Prev/Curr.isOnline

	// Construct request
	req, err := nc.NodeVerifier.SecurityVerifier.ReqFromBytes(&peerID, &data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}

	// Pass to cache req handler
	reqS, err := nc.handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}

	return reqS, nil
}

// Automatically format request and process it
func (nc *NodeCache) handleRequest(req types.Request) (interface{}, error) {
	reqPrefix, err := nc.SecurityVerifier.DetectReqPrefix(req.Data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Err while verifying req prefix: %v", err)
	}

	reqDataStr, err := nc.VerifyReqTypes(reqPrefix, req.Data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] Coulnd't verify req types. Err: %v", err)
	}

	// Create a struct based on request prefix, data
	reqS, err := nc.AutoConstructRequest(reqPrefix, reqDataStr)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Couldn't construct req. Err: %v", err)
	}

	switch reqPrefix {
	case security.MOVE_PREFIX:
		err := nc.updatePlayerPos(req.PeerID, reqS.(*types.PlayerPosition))
		return reqS, err
	case security.CHAT_PREFIX:
		return reqS, nil
	default:
		return nil, fmt.Errorf("[NODE CACHE] - Unknown request (unsupported prefix). %s", "")
	}
}

func (nc *NodeCache) updatePlayerPos(peerID string, pp *types.PlayerPosition) error {
	// Init cache map if first request
	if len(nc.Memory) == 0 {
		nc.Memory = make(Memory)
	}

	// Init playerState if needed
	if nc.Memory[peerID].Prev == nil || nc.Memory[peerID].Curr == nil {
		var states PlayerStates
		states.Prev.Position, states.Curr.Position = pp.Position, pp.Position
		nc.Memory[peerID] = states
		debug.DebugLog(fmt.Sprintf("[NODE CACHE] - First pkg from peerID: %s", peerID), debug.PLAYER)
		return nil
	}

	// Verify velocity
	currPos := nc.Memory[peerID].Curr.PlayerPosition
	validVel := nc.EngineVerifier.VerifyMoveVelocity(currPos, pp)
	if !validVel {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position. %s", "")
	}

	// Verify move position sequence
	validMove := nc.EngineVerifier.VerifyMoveDirection(currPos.DestPosition, pp.Position)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move. %s", "")
	}

	// Update playerStates with new data (move pStates stack)
	nc.Memory[peerID].Prev.Position, nc.Memory[peerID].Curr.Position = (*types.Position)(currPos.DestPosition), pp.Position
	return nil
}
