package nodes

import (
	"ticoma/packages/network/nodes/core"
)

// Basic GameNode, integral part of PlayerNode
type GameNode struct {
	*core.GameNodeCore
}

func NewGameNode() *GameNode {
	return &GameNode{
		GameNodeCore: &core.GameNodeCore{},
	}
}
