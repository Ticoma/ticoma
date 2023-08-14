package gamemap

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draw main game map in specified pos
//
// Params:
//
// startX, startY - left-most/top-most X,Y points, starting point of drawing
//
// blockSize - size of one block in px
//
// mapSize - size of entire map (in blocks)
func DrawGameMap(startX int, startY int, blockSize int, mapSizeInBlocks int) {
	mapSizeInPx := mapSizeInBlocks * blockSize
	for i := 0; i <= mapSizeInBlocks; i++ {
		rl.DrawLine(int32(startX+i*blockSize), int32(startY), int32(startX+i*blockSize), int32(mapSizeInPx)+int32(startY), rl.Black)
		for j := 0; j <= mapSizeInBlocks; j++ {
			rl.DrawLine(int32(startX), int32(startY+j*blockSize), int32(startX+mapSizeInPx), int32(startY+j*blockSize), rl.Black)
		}
	}
}
