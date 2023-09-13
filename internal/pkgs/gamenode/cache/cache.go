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
	prev *types.Player // Prev isOnline => Did such peerID ever play the game? (Login checker)
	curr *types.Player // Curr isOnline => Is this player currently online (Actual state)
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

// Did such peerID ever create an account / play?
func (nc *NodeCache) playerExists(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].prev.IsOnline
	} else {
		return false
	}
}

// Is peerID currently online?
func (nc *NodeCache) playerOnline(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].curr.IsOnline
	} else {
		return false
	}
}

// Get full playerState
func (nc *NodeCache) GetPlayer(peerID string) PlayerStates {
	if nc.playerOnline(peerID) {
		return nc.Memory[peerID]
	} else {
		return PlayerStates{}
	}
}

func (nc *NodeCache) GetPrevPlayerPos(peerID string) types.Position {
	if nc.playerOnline(peerID) {
		return *nc.Memory[peerID].prev.Position
	} else {
		return types.Position{}
	}
}

func (nc *NodeCache) GetCurrPlayerPos(peerID string) types.Position {
	if nc.playerOnline(peerID) {
		return *nc.Memory[peerID].curr.Position
	} else {
		return types.Position{}
	}
}

// Put new data to NodeCache
//
// Returns automatically constructed request interface (based on data preifx)
func (nc *NodeCache) Put(peerID string, data []byte) (interface{}, error) {

	// Construct request
	req, err := nc.NodeVerifier.SecurityVerifier.ReqFromBytes(&peerID, &data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}
	debug.DebugLog(fmt.Sprintf("[CACHE] - Request constructed. Data: {peerID: %s, data: %s}", req.PeerID, string(req.Data)), debug.PLAYER)

	// Pass to cache req handler
	reqS, err := nc.handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Req not accepted. Err: %v", err)
	}
	debug.DebugLog(fmt.Sprintf("[CACHE] - Request data stringified: %s", reqS), debug.PLAYER)

	return reqS, nil
}

// Main request handler and sorter
func (nc *NodeCache) handleRequest(req types.Request) (interface{}, error) {
	reqPrefix, err := nc.SecurityVerifier.DetectReqPrefix(req.Data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Err while verifying req prefix: %v", err)
	}
	debug.DebugLog(fmt.Sprintf("[CACHE] - Request prefix detected: %s", reqPrefix), debug.PLAYER)

	reqDataStr, err := nc.VerifyReqTypes(reqPrefix, req.Data)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] Coulnd't verify req types. Err: %v", err)
	}
	debug.DebugLog("[CACHE] - Request types verified.", debug.PLAYER)

	// Create a struct based on request prefix, data
	reqS, err := nc.AutoConstructRequest(reqPrefix, reqDataStr)
	if err != nil {
		return nil, fmt.Errorf("[NODE CACHE] - Couldn't construct req. Err: %v", err)
	}

	switch reqPrefix {
	case security.MOVE_PREFIX:
		err := nc.updatePlayerPos(req.PeerID, reqS.(types.PlayerPosition))
		return reqS, err
	case security.CHAT_PREFIX:
		return reqS, nil
	case security.REGISTER_PREFIX, security.DELETE_ACC_PREFIX, security.LOGIN_PREFIX, security.LOGOUT_PREFIX:
		err := nc.handleAccountRequest(req.PeerID, reqPrefix)
		return reqS, err
	default:
		return nil, fmt.Errorf("[NODE CACHE] - Unknown request (unsupported prefix). %s", "")
	}
}

// Sub-handler for all account-related requests
func (nc *NodeCache) handleAccountRequest(peerID string, reqPrefix string) error {
	switch reqPrefix {
	case security.REGISTER_PREFIX:
		return nc.registerPlayer(peerID)
	case security.DELETE_ACC_PREFIX:
		return nc.deletePlayer(peerID)
	case security.LOGIN_PREFIX:
		return nc.loginPlayer(peerID)
	case security.LOGOUT_PREFIX:
		return nc.logoutPlayer(peerID)
	default:
		return fmt.Errorf("[NODE CACHE] - Unknown account related request.")
	}
}

func (nc *NodeCache) loginPlayer(peerID string) error {
	if nc.playerOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Player is already logged in.")
	} else {
		nc.Memory[peerID].curr.IsOnline = true
		return nil
	}
}

func (nc *NodeCache) logoutPlayer(peerID string) error {
	if nc.playerOnline(peerID) {
		nc.Memory[peerID].curr.IsOnline = false
		return nil
	} else {
		return fmt.Errorf("[NODE CACHE] - Failed to logout. Player is already logged out.")
	}
}

func (nc *NodeCache) registerPlayer(peerID string) error {
	if nc.playerExists(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to register player. Already exists.")
	} else {
		nc.Memory[peerID].prev.IsOnline = true
		return nil
	}
}

func (nc *NodeCache) deletePlayer(peerID string) error {
	if nc.playerExists(peerID) {
		nc.Memory[peerID] = PlayerStates{} // Clear
		return nil
	} else {
		return fmt.Errorf("[NODE CACHE] - Failed to delete player. Account doesn't exist.")
	}
}

func (nc *NodeCache) updatePlayerPos(peerID string, pp types.PlayerPosition) error {

	// Ignore if affected player is offline
	if !nc.playerOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Can't update pos. Player is offline!")
	}

	// // Init cache map if first request
	// if len(nc.Memory) == 0 {
	// 	nc.Memory = make(Memory)
	// }

	// Verify velocity
	currPos := nc.Memory[peerID].curr.PlayerPosition
	validVel := nc.EngineVerifier.VerifyMoveVelocity(currPos, &pp)
	if !validVel {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position. %s", "")
	}

	// Verify move position sequence
	validMove := nc.EngineVerifier.VerifyMoveDirection(currPos.DestPosition, pp.Position)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move. %s", "")
	}

	// Update playerStates with new data (move pStates stack)
	nc.Memory[peerID].prev.Position, nc.Memory[peerID].curr.Position = (*types.Position)(currPos.DestPosition), pp.Position
	return nil
}
