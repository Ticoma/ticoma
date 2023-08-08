package gamenode

import (
	relay "ticoma/internal/packages/network/gamenode/relay"
)

// Extended game node,
// Consists of: GameNodeCore + NodeRelay
type StandaloneGameNode struct {
	*IntegralGameNode
	*relay.GameNodeRelay
}

func NewStandaloneGameNode() *StandaloneGameNode {
	ign := NewIntegralGameNode()
	return &StandaloneGameNode{
		IntegralGameNode: ign,
		GameNodeRelay:    &relay.GameNodeRelay{},
	}
}
