package verifier

import (
	"ticoma/internal/packages/gamenode/cache/verifier/engine"
	"ticoma/internal/packages/gamenode/cache/verifier/security"
)

type NodeVerifier struct {
	*engine.EngineVerifier
	*security.SecurityVerifier
}

func New() *NodeVerifier {
	return &NodeVerifier{
		EngineVerifier:   &engine.EngineVerifier{},
		SecurityVerifier: &security.SecurityVerifier{},
	}
}
