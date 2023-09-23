package cache

import (
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
	"time"
)

//
// All account related request functions, handlers for Cache
//

func (nc *NodeCache) accountExists(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].Prev.IsOnline
	}
	return false
}

func (nc *NodeCache) accountOnline(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].Curr.IsOnline
	}
	return false
}

func (nc *NodeCache) handleAccountRequest(peerID string, reqPrefix string, reqS *interface{}) error {
	switch reqPrefix {
	case security.REGISTER_PREFIX:
		return nc.registerPlayer(peerID, (*reqS).(string))
	case security.DELETE_ACC_PREFIX:
		return nc.deletePlayer(peerID)
	case security.LOGIN_PREFIX:
		return nc.loginPlayer(peerID)
	case security.LOGOUT_PREFIX:
		fmt.Println("Logout request")
		return nc.logoutPlayer(peerID)
	default:
		return fmt.Errorf("[NODE CACHE] - Unknown account related request.")
	}
}

func (nc *NodeCache) registerPlayer(peerID string, nickname string) error {
	if nc.accountExists(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to register player. Already exists.")
	}
	p := nc.GetPlayer(peerID)
	ts := time.Now().UnixMilli()
	spawnPos := types.PlayerPosition{
		Timestamp:    ts,
		Position:     types.Position{X: SPAWN_POS_X, Y: SPAWN_POS_Y},
		DestPosition: types.DestPosition{X: SPAWN_POS_X, Y: SPAWN_POS_Y},
	}
	p.Prev.IsOnline, p.Curr.IsOnline = true, true
	p.Prev.PlayerPosition, p.Curr.PlayerPosition = spawnPos, spawnPos
	p.Prev.Nick, p.Curr.Nick = nickname, nickname
	nc.Memory[peerID] = *p
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - PeerID: %s created an account", peerID), debug.PLAYER)
	return nil
}

func (nc *NodeCache) deletePlayer(peerID string) error {
	if !nc.accountExists(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to delete player. Account doesn't exist.")
	}
	nc.Memory[peerID] = PlayerStates{}
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s deleted their account", peerID), debug.PLAYER)
	return nil
}

func (nc *NodeCache) loginPlayer(peerID string) error {
	if nc.accountOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Player is already logged in.")
	}
	p := nc.Memory[peerID]
	p.Curr.IsOnline = true
	nc.Memory[peerID] = p
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s logged in", peerID), debug.PLAYER)
	return nil
}

func (nc *NodeCache) logoutPlayer(peerID string) error {
	if !nc.accountOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to logout. Player is already logged out")
	}
	fmt.Println("Logging out player: ", peerID)
	p := nc.Memory[peerID]
	p.Curr.IsOnline = false
	nc.Memory[peerID] = p
	debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s logged out", peerID), debug.PLAYER)
	return nil
}
