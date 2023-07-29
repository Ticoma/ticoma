// interfaces for Nodes
// e.g. the default, go-to ActionDataPackage (ADP) / ActionDataPackageTimestamped
package interfaces

type Position struct {
	X int
	Y int
}

type DestPosition struct {
	*Position
}

type ActionDataPackage struct {
	PlayerId     int
	PlayerPubKey string
	*Position
	*DestPosition
}

type ActionDataPackageTimestamped struct {
	*ActionDataPackage
	Timestamp int
}
