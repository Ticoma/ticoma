package gamenode

import (
	core "ticoma/packages/network/gamenode/core"
	relay "ticoma/packages/network/gamenode/relay"
)

// Extended game node,
// Consists of: GameNodeCore + NodeRelay
type StandaloneGameNode struct {
	*GameNode
	*relay.GameNodeRelay
}

func NewStandaloneGameNode() *StandaloneGameNode {
	return &StandaloneGameNode{
		GameNode: &GameNode{
			GameNodeCore: &core.GameNodeCore{},
		},
		GameNodeRelay: &relay.GameNodeRelay{},
	}
}
