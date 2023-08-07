package verifier

type NodeVerifier struct {
	*EngineVerifier
	*SecurityVerifier
}

func NewNodeVerifier() *NodeVerifier {
	return &NodeVerifier{
		EngineVerifier:   &EngineVerifier{},
		SecurityVerifier: &SecurityVerifier{},
	}
}
