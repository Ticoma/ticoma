package verifier

import (
	"fmt"
	"reflect"
	. "ticoma/packages/nodes/interfaces"
)

// SecurityVerifier
type SecurityVerifier struct{}

// Is this really needed?
func (sv *SecurityVerifier) VerifyADPTypes(pkg *ActionDataPackageTimestamped) {
	t := reflect.TypeOf(pkg)
	fmt.Println("verified:", t)
}
