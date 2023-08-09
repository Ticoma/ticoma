package verifier

type NodeVerifier struct {
	*EngineVerifier
	*SecurityVerifier
}

func New() *NodeVerifier {
	return &NodeVerifier{
		EngineVerifier:   &EngineVerifier{},
		SecurityVerifier: &SecurityVerifier{},
	}
}
