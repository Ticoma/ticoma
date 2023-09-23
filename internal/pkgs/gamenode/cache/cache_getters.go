package cache

import (
	"ticoma/types"
)

func (nc *NodeCache) GetAll() *Memory {
	return &nc.Memory
}

func (nc *NodeCache) GetPlayer(peerID string) *PlayerStates {
	if !nc.accountExists(peerID) {
		return &PlayerStates{}
	}
	p := nc.Memory[peerID]
	return &p
}

func (nc *NodeCache) GetPrevPlayerPos(peerID string) *types.PlayerPosition {
	if !nc.accountOnline(peerID) {
		return &types.PlayerPosition{}
	}
	p := nc.Memory[peerID]
	return &p.Prev.PlayerPosition
}

func (nc *NodeCache) GetCurrPlayerPos(peerID string) *types.PlayerPosition {
	if !nc.accountOnline(peerID) {
		return &types.PlayerPosition{}
	}
	p := nc.Memory[peerID]
	return &p.Curr.PlayerPosition
}

func (nc *NodeCache) GetNickname(peerID string) *string {
	if !nc.accountExists(peerID) {
		return new(string)
	}
	p := nc.Memory[peerID]
	return &p.Curr.Nick
}
