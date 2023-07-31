package nodes

import (
	"fmt"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules"
	. "ticoma/packages/nodes/modules/verifier"
)

func PlayerNode() {
	fmt.Printf("Hello from PlayerNode\n")

	// Init two players
	v := NewVerifier()
	nc := NewNodeCache(v)

	// init player 0
	pkg := ActionDataPackageTimestamped{
		ActionDataPackage: &ActionDataPackage{
			PlayerId:     0,
			PubKey:       "PUBKEY",
			Position:     &Position{X: 1, Y: 1},
			DestPosition: &DestPosition{X: 1, Y: 1},
		},
		Timestamp: 1,
	}

	nc.Put(pkg)

	// init second node later
	v2 := NewVerifier()
	nc2 := NewNodeCache(v2)

	cache := nc2.GetAll()

	fmt.Println(cache)

}
