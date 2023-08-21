package gamemap

import "ticoma/client/packages/utils"

type GameMap struct {
	FrontLayer [][]MapBlock
	BgLayer    [][]MapBlock
}

type MapBlock struct {
	BlockId     int
	Collision   bool
	Interactive bool
}

// Empty map constructor
func NewMap(width int, height int) *GameMap {
	layer := make([][]MapBlock, height)
	for i := range layer {
		layer[i] = make([]MapBlock, width)
	}
	return &GameMap{
		FrontLayer: layer,
		BgLayer:    layer,
	}
}

// Fill map layer with specified layer config/configs
func (gm *GameMap) LoadMap(frontLayer [][]MapBlock, bgLayer [][]MapBlock) {
	gm.FrontLayer = frontLayer
	gm.BgLayer = bgLayer
}

// Tmp: generates a mosaic map full of random blocks
func (gm *GameMap) GenerateRandomMap() {
	for i := 0; i < len(gm.BgLayer); i++ {
		for j := 0; j < len(gm.BgLayer[0]); j++ {
			randBlockId := utils.RandRange(0, 2)
			gm.BgLayer[i][j] = MapBlock{
				BlockId:     randBlockId,
				Collision:   false,
				Interactive: false,
			}
		}
	}
}

// Insert block on {x, y} pos on specified map layer
func (gm *GameMap) InsertBlock(gamemap *GameMap, newBlock MapBlock, x int, y int, updateOnBgLayer bool) {
	if updateOnBgLayer {
		gamemap.BgLayer[x][y] = newBlock
	}
	gamemap.FrontLayer[x][y] = newBlock
}
