package verifier

type Verifier struct {
	*EngineVerifier
	*SecurityVerifier
}

func NewVerifier() *Verifier {
	return &Verifier{
		EngineVerifier:   &EngineVerifier{},
		SecurityVerifier: &SecurityVerifier{},
	}
}
