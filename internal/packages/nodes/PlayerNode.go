package nodes

import (
	"fmt"
	// . "ticoma/packages/nodes/interfaces"
	// . "ticoma/packages/nodes/modules/verifier"
)

func PlayerNode() {
	fmt.Printf("Hello from PlayerNode\n")

	// testADPT := ActionDataPackageTimestamped{
	// 	&ActionDataPackage{
	// 		PlayerId:     1,
	// 		PlayerPubKey: "PUBKEY",
	// 		Position:     &Position{X: 1, Y: 1},
	// 		DestPosition: &DestPosition{&Position{X: 1, Y: 1}},
	// 	},
	// 	123, // Timestamp
	// }

	// testVerifier := NewVerifier()
	// testVerifier.VerifyPlayerMovement()
	// testVerifier.VerifyADPTypes()

	// pos := Position{X: 1, Y: 1}
	// fmt.Printf("Pos: %v\n", pos)
}

func Add(a, b int) int {
	return a + b
}
