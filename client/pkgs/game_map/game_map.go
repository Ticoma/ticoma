package gamemap

//
// Needs a rework, (soon)
//

type GameMap struct {
	FrontLayer [][]MapBlock
	BgLayer    [][]MapBlock
}

type MapBlock struct {
	BlockId     int
	Collision   bool
	Interactive bool
}
