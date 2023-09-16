package cache

import (
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache/verifier"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
	"time"
)

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

func (nc *NodeCache) GetAll() Memory {
	return nc.Memory
}

// Did such peerID ever create an account / play?
func (nc *NodeCache) playerExists(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		debug.DebugLog("Player exists", debug.PLAYER)
		return nc.Memory[peerID].Prev.IsOnline
	} else {
		return false
	}
}

// Is peerID Currently online?
func (nc *NodeCache) playerOnline(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].Curr.IsOnline
	} else {
		return false
	}
}

// Returns ptr to existing player in cache (if doesn't exist -> empty playerState type)
func (nc *NodeCache) GetPlayer(peerID string) *PlayerStates {
	if nc.playerExists(peerID) {
		p := nc.Memory[peerID]
		return &p
	} else {
		return &PlayerStates{}
	}
}

func (nc *NodeCache) GetPrevPlayerPos(peerID string) types.Position {
	if nc.playerOnline(peerID) {
		return nc.Memory[peerID].Prev.Position
	} else {
		return types.Position{}
	}
}

func (nc *NodeCache) GetCurrPlayerPos(peerID string) types.Position {
	if nc.playerOnline(peerID) {
		return nc.Memory[peerID].Curr.Position
	} else {
		return types.Position{}
	}
}

// Auto detect request type and try to put new data in Cache
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

	debug.DebugLog(fmt.Sprintf("[CACHE] - Request handled"), debug.PLAYER)
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
	debug.DebugLog("[CACHE] - Request constructed successfully - passing to handler.", debug.PLAYER)

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

// Sub-handler for all account related requests
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
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s Login request", peerID), debug.PLAYER)
	if nc.playerOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Player is already logged in.")
	} else {
		p := nc.Memory[peerID]
		p.Curr.IsOnline = true
		nc.Memory[peerID] = p
		return nil
	}
}

func (nc *NodeCache) logoutPlayer(peerID string) error {
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s Logout request", peerID), debug.PLAYER)
	if nc.playerOnline(peerID) {
		p := nc.Memory[peerID]
		p.Curr.IsOnline = false
		nc.Memory[peerID] = p
		return nil
	} else {
		return fmt.Errorf("[NODE CACHE] - Failed to logout. Player is already logged out.")
	}
}

func (nc *NodeCache) registerPlayer(peerID string) error {
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s Register request", peerID), debug.PLAYER)
	if nc.playerExists(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to register player. Already exists.")
	} else {
		p := nc.GetPlayer(peerID)
		// Online state
		p.Prev.IsOnline = true
		p.Curr.IsOnline = true
		// Init spawn position
		ts := time.Now().UnixMilli()
		debug.DebugLog(fmt.Sprintf("TIMESTAMP REGISTER: %d", ts), debug.PLAYER)
		spawnPos := types.PlayerPosition{
			Timestamp:    ts,
			Position:     types.Position{X: 13, Y: 13},     // SPAWN POS (CHANGE LATER)
			DestPosition: types.DestPosition{X: 13, Y: 13}, // SPAWN POS (CHANGE LATER)
		}
		p.Prev.PlayerPosition = spawnPos
		p.Curr.PlayerPosition = spawnPos
		nc.Memory[peerID] = *p
		return nil
	}
}

func (nc *NodeCache) deletePlayer(peerID string) error {
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s Delete acc request", peerID), debug.PLAYER)
	if nc.playerExists(peerID) {
		nc.Memory[peerID] = PlayerStates{}
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

	p := nc.Memory[peerID]

	// Verify velocity
	currPos := p.Curr.PlayerPosition
	validVel := nc.EngineVerifier.VerifyMoveVelocity(&currPos, &pp)
	if !validVel {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position. %s", "")
	}

	// Verify move position sequence
	validMove := nc.EngineVerifier.VerifyMoveDirection(&currPos.DestPosition, &pp.Position)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move. %s", "")
	}

	// Update playerStates with new data (move pStates stack)
	p.Prev.PlayerGameData.Position = types.Position(p.Prev.PlayerGameData.DestPosition)
	p.Prev.DestPosition = types.DestPosition(p.Curr.Position)
	p.Curr.Position = pp.Position
	p.Curr.DestPosition = pp.DestPosition
	nc.Memory[peerID] = p

	return nil
}
