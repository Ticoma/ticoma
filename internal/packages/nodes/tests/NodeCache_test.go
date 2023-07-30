package tests

import (
	"fmt"
	"testing"

	// . "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules"
	. "ticoma/packages/nodes/modules/verifier"
)

func TestBasicNodeCache(t *testing.T) {

	// init
	v := NewVerifier()
	nc := NewNodeCache(v)

	// no player, solo cache
	cache := nc.GetCache(0)
	fmt.Println(cache)

	// pkg1 := ActionDataPackageTimestamped{
	// 	ActionDataPackage: &ActionDataPackage{
	// 		PlayerId:     1,
	// 		PlayerPubKey: "PUBKEY",
	// 		Position:     &Position{X: 1, Y: 1},
	// 		DestPosition: &DestPosition{Position: &Position{X: 2, Y: 2}},
	// 	},
	// 	Timestamp: 1000,
	// }

	// pkg2 := ActionDataPackageTimestamped{
	// 	ActionDataPackage: &ActionDataPackage{
	// 		PlayerId:     1,
	// 		PlayerPubKey: "PUBKEY",
	// 		Position:     &Position{X: 2, Y: 2},
	// 		DestPosition: &DestPosition{Position: &Position{X: 2, Y: 2}},
	// 	},
	// 	Timestamp: 2000,
	// }

	// got := v.EngineVerifier.VerifyPlayerMovement(&pkg1, &pkg2)
	// want := true

	// if got != want {
	// 	t.Errorf("got %t, wanted %t", got, want)
	// }

}
