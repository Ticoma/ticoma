package nodes

import (
	"ticoma/packages/network/nodes/relay"
)

// Extended game node,
// Consists of: GameNodeCore + NodeRelay
type StandaloneGameNode struct {
	*GameNode
	*relay.GameNodeRelay
}
