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

// Did such peerID ever create an account / play?
func (nc *NodeCache) accountExists(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].Prev.IsOnline
	} else {
		return false
	}
}

// Is peerID currently online?
func (nc *NodeCache) accountOnline(peerID string) bool {
	_, exists := nc.Memory[peerID]
	if exists {
		return nc.Memory[peerID].Curr.IsOnline
	} else {
		return false
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

// Create a new account
func (nc *NodeCache) registerPlayer(peerID string) error {
	if nc.accountExists(peerID) {
		return fmt.Errorf("[NODE CACHE] - Failed to register player. Already exists.")
	} else {
		// Create a playerState @ spawnPos with current timestamp, set status to online
		p := nc.GetPlayer(peerID)
		p.Prev.IsOnline, p.Curr.IsOnline = true, true
		ts := time.Now().UnixMilli()
		spawnPos := types.PlayerPosition{
			Timestamp:    ts,
			Position:     types.Position{X: SPAWN_POS_X, Y: SPAWN_POS_Y},
			DestPosition: types.DestPosition{X: SPAWN_POS_X, Y: SPAWN_POS_Y},
		}
		p.Prev.PlayerPosition, p.Curr.PlayerPosition = spawnPos, spawnPos
		nc.Memory[peerID] = *p
		debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s created an account", peerID), debug.PLAYER)
		return nil
	}
}

// Delete an existing account
func (nc *NodeCache) deletePlayer(peerID string) error {
	if nc.accountExists(peerID) {
		nc.Memory[peerID] = PlayerStates{}
		debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s deleted their account", peerID), debug.PLAYER)
		return nil
	} else {
		return fmt.Errorf("[NODE CACHE] - Failed to delete player. Account doesn't exist.")
	}
}

// Login to an existing account
func (nc *NodeCache) loginPlayer(peerID string) error {
	if nc.accountOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Player is already logged in.")
	} else {
		p := nc.Memory[peerID]
		p.Curr.IsOnline = true
		nc.Memory[peerID] = p
		debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s logged in", peerID), debug.PLAYER)
		return nil
	}
}

// Logout from account
func (nc *NodeCache) logoutPlayer(peerID string) error {
	if nc.accountOnline(peerID) {
		p := nc.Memory[peerID]
		p.Curr.IsOnline = false
		nc.Memory[peerID] = p
		debug.DebugLog(fmt.Sprintf("[NODE CACHE] - Player %s logged out", peerID), debug.PLAYER)
		return nil
	} else {
		return fmt.Errorf("[NODE CACHE] - Failed to logout. Player is already logged out")
	}
}
