package nodes

import (
	"fmt"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules"
	. "ticoma/packages/nodes/modules/verifier"
)

func PlayerNode() {
	fmt.Printf("Hello from PlayerNode\n")

	v := NewVerifier()
	nc := NewNodeCache(v)

	pkg := ActionDataPackageTimestamped{
		ActionDataPackage: &ActionDataPackage{
			PlayerId:     1,
			PubKey:       "PUBKEY",
			Position:     &Position{X: 1, Y: 1},
			DestPosition: &DestPosition{X: 2, Y: 2},
		},
		Timestamp: 1337,
	}

	// should return empty map
	// fmt.Println("EMPTY: ", nc.GetAll())

	nc.Put(pkg)

	// fmt.Println(nc.GetCache(1))

	// init for player 2
	// v2 := NewVerifier()
	// nc2 := NewNodeCache(v2)

	pkg2 := ActionDataPackageTimestamped{
		ActionDataPackage: &ActionDataPackage{
			PlayerId:     2,
			PubKey:       "PUBKEY2",
			Position:     &Position{X: 3, Y: 3},
			DestPosition: &DestPosition{X: 3, Y: 3},
		},
		Timestamp: 7331,
	}

	nc.Put(pkg2)

	fmt.Println("PLAYER 1's CACHE: ")
	fmt.Println(nc.GetCache(1))
	fmt.Println(nc.GetCache(2))

	// check
	pkg3 := ActionDataPackageTimestamped{
		ActionDataPackage: &ActionDataPackage{
			PlayerId:     1,
			PubKey:       "PUBKEY",
			Position:     &Position{X: 2, Y: 2},
			DestPosition: &DestPosition{X: 2, Y: 2},
		},
		Timestamp: 2337,
	}

	nc.Put(pkg3)

	fmt.Println("PLAYER 1's CACHE: ")
	fmt.Println(nc.GetCache(1))
	fmt.Println(nc.GetCache(2))
}
