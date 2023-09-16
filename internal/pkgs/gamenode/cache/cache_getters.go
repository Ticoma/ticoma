package cache

import (
	"ticoma/types"
)

//
// All public getters for Cache
//

// Get entire cache
func (nc *NodeCache) GetAll() *Memory {
	return &nc.Memory
}

// Get ptr to existing player in cache (if doesn't exist -> empty playerState type)
func (nc *NodeCache) GetPlayer(peerID string) *PlayerStates {
	if nc.accountExists(peerID) {
		p := nc.Memory[peerID]
		return &p
	} else {
		return &PlayerStates{}
	}
}

// Get previous player position (must be online)
func (nc *NodeCache) GetPrevPlayerPos(peerID string) *types.PlayerPosition {
	if nc.accountOnline(peerID) {
		p := nc.Memory[peerID]
		return &p.Prev.PlayerPosition
	} else {
		return &types.PlayerPosition{}
	}
}

// Get previous player position (must be online)
func (nc *NodeCache) GetCurrPlayerPos(peerID string) *types.PlayerPosition {
	if nc.accountOnline(peerID) {
		p := nc.Memory[peerID]
		return &p.Curr.PlayerPosition
	} else {
		return &types.PlayerPosition{}
	}
}
