package cache

import (
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache/verifier"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
)

const SPAWN_POS_X = 6
const SPAWN_POS_Y = 6

type NodeCache struct {
	Memory
	*verifier.NodeVerifier
}

type PlayerStates struct {
	Prev types.Player // Prev isOnline => Did such peerID ever play the game? (Login checker)
	Curr types.Player // Curr isOnline => Is this player Currently online (Actual state)
}

type Memory map[string]PlayerStates

func New() *NodeCache {
	v := verifier.New()
	return &NodeCache{
		Memory:       make(Memory),
		NodeVerifier: v,
	}
}

// Auto detect request type and try to put new data in Cache
//
// Returns constructed request and its prefix
func (nc *NodeCache) Put(peerID string, data []byte) (interface{}, string, error) {

	// Construct request
	req, err := nc.NodeVerifier.SecurityVerifier.ReqFromBytes(&peerID, &data)
	if err != nil {
		return nil, "", fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}
	debug.DebugLog(fmt.Sprintf("[CACHE] - Request constructed. Req: {peerID: %s, data: \"%s\"}", req.PeerID, string(req.Data)), debug.PLAYER)

	// Pass to cache req handler
	reqS, reqPfx, err := nc.handleRequest(req)
	if err != nil {
		return nil, "", fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}

	debug.DebugLog(fmt.Sprintf("[CACHE] - Request handled"), debug.PLAYER)
	return *reqS, reqPfx, nil
}

// Main request handler and sorter
func (nc *NodeCache) handleRequest(req types.Request) (*interface{}, string, error) {

	// Detect prefix from request data
	reqPrefix, err := nc.SecurityVerifier.DetectReqPrefix(req.Data)
	if err != nil {
		return nil, "", fmt.Errorf("[NODE CACHE] - Err while verifying req prefix: %v", err)
	}
	debug.DebugLog(fmt.Sprintf("[CACHE] - Request prefix detected: %s", reqPrefix), debug.PLAYER)

	// Verify types with request schema
	reqDataStr, err := nc.VerifyReqTypes(reqPrefix, req.Data)
	if err != nil {
		return nil, "", fmt.Errorf("[NODE CACHE] Coulnd't verify req types. Err: %v", err)
	}
	debug.DebugLog("[CACHE] - Request types verified.", debug.PLAYER)

	// Create a struct based on request prefix and data
	reqS, err := nc.AutoConstructRequest(reqPrefix, reqDataStr, req.PeerID)
	if err != nil {
		return nil, "", fmt.Errorf("[NODE CACHE] - Couldn't construct req. Err: %v", err)
	}
	debug.DebugLog("[CACHE] - Request constructed successfully - passing to handler.", debug.PLAYER)

	switch reqPrefix {
	// Game request
	case security.MOVE_PREFIX, security.CHAT_PREFIX:
		err := nc.handleGameRequest(req.PeerID, reqPrefix, &reqS)
		return &reqS, reqPrefix, err
	// Account request
	case security.REGISTER_PREFIX, security.DELETE_ACC_PREFIX, security.LOGIN_PREFIX, security.LOGOUT_PREFIX:
		err := nc.handleAccountRequest(req.PeerID, reqPrefix, &reqS)
		return &reqS, reqPrefix, err
	// Unknown
	default:
		return nil, "", fmt.Errorf("[NODE CACHE] - Unknown request (unsupported prefix). %s", "")
	}
}
