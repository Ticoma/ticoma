package gamenode

import (
	core "ticoma/internal/packages/network/gamenode/core"
	relay "ticoma/internal/packages/network/gamenode/relay"
)

// Extended game node,
// Consists of: GameNodeCore + NodeRelay
type StandaloneGameNode struct {
	*IntegralGameNode
	*relay.GameNodeRelay
}

func NewStandaloneGameNode() *StandaloneGameNode {
	return &StandaloneGameNode{
		IntegralGameNode: &IntegralGameNode{
			GameNodeCore: &core.GameNodeCore{},
		},
		GameNodeRelay: &relay.GameNodeRelay{},
	}
}
