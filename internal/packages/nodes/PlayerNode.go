package nodes

import (
	"encoding/json"
	"fmt"

	// . "ticoma/packages/nodes/interfaces"
	// . "ticoma/packages/nodes/modules"
	. "ticoma/packages/nodes/modules/verifier"
)

func PlayerNode() {
	fmt.Printf("Hello from PlayerNode\n")

	// testADPT := ActionDataPackageTimestamped{
	// 	&ActionDataPackage{
	// 		PlayerId:     1,
	// 		PubKey: "PUBKEY",
	// 		Position:     &Position{X: 1, Y: 1},
	// 		DestPosition: &DestPosition{&Position{X: 1, Y: 1}},
	// 	},
	// 	123, // Timestamp
	// }

	// testADP := ActionDataPackage{
	// 	PlayerId:     1,
	// 	PlayerPubKey: "PUBKEY",
	// 	Position:     &Position{X: 1, Y: 1},
	// 	DestPosition: &DestPosition{X: 3, Y: 3},
	// }

	// // test := InterfaceToString(testADP)
	// test, err := json.Marshal(testADP)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// testADP_str := fmt.Sprint(testADP)
	// fmt.Println(testADP_str)

	testVerifier := NewVerifier()

	// incorrect := []byte(`{"PLAYER_ID":"test","pubKey":"PUBKEY","zzz":{"posX":1,"posY":1},"destPos":{"destPosX":1,"destPosY":1}}`)
	correct := []byte(`{"playerId":0,"pubKey":"PUBKEY","pos":{"posX":1,"posY":1},"destPos":{"destPosX":1,"destPosY":1}}`)

	res := testVerifier.VerifyADPTypes(correct)
	fmt.Println(res)

	// pos := Position{X: 1, Y: 1}
	// fmt.Printf("Pos: %v\n", pos)

	// v := NewVerifier()
	// cache := NewNodeCache(v)
	// test := cache.GetCache(0)
	// fmt.Println(test)

}

func InterfaceToString(obj interface{}) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}
