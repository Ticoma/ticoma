// interfaces for Nodes
// e.g. the default, go-to ActionDataPackage (ADP) / ActionDataPackageTimestamped
package interfaces

type Position struct {
	X int `json:"posX"`
	Y int `json:"posY"`
}

type DestPosition struct {
	X int `json:"destPosX"`
	Y int `json:"destPosY"`
}

type ActionDataPackage struct {
	PlayerId      int    `json:"playerId"`
	PubKey        string `json:"pubKey"`
	*Position     `json:"pos"`
	*DestPosition `json:"destPos"`
}

type ActionDataPackageTimestamped struct {
	*ActionDataPackage
	Timestamp int `json:"timestamp"`
}
