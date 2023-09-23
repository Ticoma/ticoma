package gamemap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	c "ticoma/client/pkgs/constants"

	"github.com/solarlune/paths"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameMap struct {
	Initialized bool
	Txt         rl.RenderTexture2D
	PathMap     *paths.Grid
}

func New() *GameMap {
	return &GameMap{
		Initialized: false,
		Txt:         rl.RenderTexture2D{},
		PathMap:     nil,
	}
}

// Load map from file, verify contents and draw terrain
func (gm *GameMap) Init(path string, spriteSource *rl.Texture2D) error {

	// Verify & read layers, contents
	mSize, mContent, lCount, lLen, err := readMapFile(path)
	if err != nil {
		return fmt.Errorf("Couldn't load map. Err: %s", err.Error())
	}
	// Init texture for map
	mapSizePx := int32(int(c.BLOCK_SIZE) * mSize)
	gm.Txt = rl.LoadRenderTexture(mapSizePx, mapSizePx)

	// Fill map texture with blocks
	cMap, err := gm.drawMapBlocks(spriteSource, mSize, mContent, lCount, lLen)
	if err != nil {
		return fmt.Errorf("Couldn't draw map blocks. Err: %s", err.Error())
	}

	// Generate path & collision map
	err = gm.createPathMap(cMap)
	if err != nil {
		return fmt.Errorf("Couldn't create path map. Err: %s", err.Error())
	}

	gm.Initialized = true
	return nil
}

// Returns: mapSize, mapContent, layerCount, layerLength
func readMapFile(path string) (int, string, int, int, error) {
	// Read map file
	bytes, err := os.ReadFile(path)
	if err != nil {
		return 0, "", 0, 0, fmt.Errorf("Can't find map file. Check your path")
	}
	// Read map size
	mapSizeEndChar := ";"
	index := strings.Index(string(bytes), mapSizeEndChar)
	if index == -1 {
		return 0, "", 0, 0, fmt.Errorf("Couldn't read map size from file. Is the map format correct?")
	}
	// TODO: add support for asymmetrical map sizes
	mapSize, err := strconv.Atoi(string(bytes[:index]))
	if err != nil {
		return 0, "", 0, 0, fmt.Errorf("Couldn't convert map size to integer")
	}

	// fmt.Println("MAP SIZE: ", mapSize)
	layerLen := (mapSize * mapSize) * 2 // including commas
	mapContent := string(bytes[index+1:])
	layerCount := len(mapContent) / layerLen // cut in half to get real number of layers (don't count commas)

	// Check if all layers are present in file
	expectedLayersContentLen := layerCount * layerLen
	validLayersContent := expectedLayersContentLen == len(mapContent)

	if !validLayersContent {
		return 0, "", 0, 0, fmt.Errorf("Expected layers content don't match file content. Layers: %d\nExpected length of layers content: %d, Got: %d", layerCount, expectedLayersContentLen, len(mapContent))
	}

	return mapSize, mapContent, layerCount, layerLen, nil
}

// Fill game map texture with blocks
func (gm *GameMap) drawMapBlocks(spriteSource *rl.Texture2D, mSize int, mContent string, lCount int, lLen int) ([]string, error) {

	var cMap []string
	var cMapLayer string

	rl.BeginTextureMode(gm.Txt)
	for i := 0; i < lCount; i++ {

		layerContent := mContent[i*lLen : (i+1)*lLen]
		content := strings.Split(layerContent, ",")
		cMapLayer = ""

		for j := 0; j < len(content)-1; j++ {
			// Get block position based on its index in mapContent array
			blockPos := rl.Vector2{
				X: float32((j % mSize) * int(c.BLOCK_SIZE)),
				Y: float32((j / mSize) * int(c.BLOCK_SIZE)),
			}
			// Draw each block on texture
			err := gm.drawBlock(spriteSource, content[j], &blockPos)
			if err != nil {
				return nil, fmt.Errorf("Couldn't draw block on map. Err: %s", err.Error())
			}
			// Add block representation to cMap
			if i == 0 { // Only one (BG LAYER) layer for now
				switch content[j] {
				case "0": // 0 == transparent block, col on
					cMapLayer += string(c.COLLISION_BLOCK_RUNE)
				default: // other == should be walkable
					cMapLayer += string(c.WALKABLE_BLOCK_RUNE)
				}
			}
			// update cMap layer
			if len(cMapLayer) == mSize {
				cMap = append(cMap, cMapLayer)
				cMapLayer = ""
			}
		}
	}

	rl.EndTextureMode()
	return cMap, nil
}

// Create grid PathMap from string collision map
func (gm *GameMap) createPathMap(cMap []string) error {

	// fmt.Println("PATH MAP: ")
	// fmt.Println(len(cMap))
	// for _, cl := range cMap {
	// 	fmt.Printf("%s\n", cl)
	// }

	if gm.PathMap != nil {
		return fmt.Errorf("Path map already exists for this map")
	}
	gm.PathMap = paths.NewGridFromStringArrays(cMap, 1, 1)
	gm.PathMap.SetWalkable('x', false) // Collision
	return nil
}
func (gm *GameMap) drawBlock(spriteSource *rl.Texture2D, id string, blockPos *rl.Vector2) error {
	blockId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Err convert block id. Err: %s", err.Error())
	}
	block := rl.Rectangle{
		X:      float32(blockId * int(c.BLOCK_SIZE)),
		Y:      0,
		Width:  c.BLOCK_SIZE,
		Height: c.BLOCK_SIZE,
	}
	// Ignore block 0 (transparent)
	if blockId != 0 {
		rl.DrawTextureRec(*spriteSource, block, *blockPos, rl.White)
	}
	return nil
}

func (gm *GameMap) Render() {
	if !gm.Initialized {
		fmt.Println("Can't render. Map not initialized")
		return
	}
	rl.DrawTextureRec(gm.Txt.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(gm.Txt.Texture.Width), Height: float32(-gm.Txt.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
}

// Check collision for a block on a Grid. Returns value of cell.Walkable
func IsTileWalkable(grid *paths.Grid, x int, y int) bool {
	tile := grid.Get(x, y)
	if tile != nil {
		return tile.Walkable
	}
	return false
}
