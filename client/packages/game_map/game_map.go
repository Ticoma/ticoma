package gamemap

import (
	"ticoma/client/packages/utils"
)

type GameMap struct {
	FrontLayer [][]MapBlock
	BgLayer    [][]MapBlock
}

type MapBlock struct {
	BlockId     int
	Collision   bool
	Interactive bool
}

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

func (gm *GameMap) GenerateRandomMap(gameMap GameMap) {
	for i := 0; i < len(gameMap.BgLayer); i++ {
		for j := 0; j < len(gameMap.BgLayer[0]); i++ {
			randBlockId := utils.RandRange(0, 4)
			gameMap.BgLayer[i][j] = MapBlock{
				BlockId:     randBlockId,
				Collision:   false,
				Interactive: false,
			}
		}
	}
}

func (gm *GameMap) UpdateBlock(gamemap *GameMap, newBlock MapBlock, x int, y int, updateOnBgLayer bool) {
	if updateOnBgLayer {
		gamemap.BgLayer[x][y] = newBlock
	}
	gamemap.FrontLayer[x][y] = newBlock
}
