package game

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	c "ticoma/client/pkgs/constants"

	"ticoma/client/pkgs/camera"
	// "ticoma/client/pkgs/input/keyboard"
	// "ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"

	left "ticoma/client/pkgs/drawing/scenes/game/panels/left"
	right "ticoma/client/pkgs/drawing/scenes/game/panels/right"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game variables
var chatInput []byte
var chatMsgs []string
var inputHold int

var gameCam *camera.GameCamera

// Textures, assets
var gameMap rl.RenderTexture2D
var blocksTxt rl.Texture2D

// var blocksImg *rl.Image

// Side panels
var leftPanel *left.LeftPanel
var rightPanel *right.RightPanel

// Misc
var sceneReady bool = false
var SIDE_PANEL_WIDTH int32

// Load all the textures, assets needed to render game scene
func loadGameScene() {

	// Load sprites
	blocksTxt = rl.LoadTexture("../client/assets/textures/sprites/blocks.png")

	err := ImportMapFromFile("")
	if err != nil {
		fmt.Println("ETLkhdfklhgkdlgh ERR")
		return
	}

	// blocksImg = rl.LoadImage("../client/assets/textures/sprites/blocks.png")
	// blocksTxt = rl.LoadTextureFromImage(blocksImg)

	// Setup res, scaling
	// SIDE_PANEL_WIDTH = int32((c.SCREEN.Width / 4))

	// Setup game
	gameCam = camera.New(float32(gameMap.Texture.Width/2), float32(gameMap.Texture.Height/2), float32(c.SCREEN.Width/2), float32((c.SCREEN.Height)/2))

	// tmp, Draw map on world from texture
	// spawnTxt = rl.LoadTextureFromImage(spawnImg)

	// Init side panels
	// leftTabs := map[int][2]string{
	// 	0: {"Chat", "C"},
	// 	1: {"Build info", "B"},
	// 	2: {"Tabssss bro", "Tb"},
	// }
	// leftPanel = left.New(float32(SIDE_PANEL_WIDTH), float32(c.SCREEN.Height), 0, 0, &c.COLOR_PANEL_BG, leftTabs)

	// rightTabs := map[int][2]string{
	// 	0: {"Inventory", "I"},
	// 	1: {"Settings", "S"},
	// }
	// rightPanel = right.New(float32(SIDE_PANEL_WIDTH), float32(c.SCREEN.Height), float32(int32(c.SCREEN.Width)-SIDE_PANEL_WIDTH), 0, &c.COLOR_PANEL_BG, rightTabs)

	sceneReady = true
}

// Render game scene (requires Player to be logged in)
func RenderGameScene(cp *player.ClientPlayer) {

	if !sceneReady {
		loadGameScene()
	}

	// // Draw players
	// // DrawMap(&world, &spawnTxt, gameCam.Zoom)
	// DrawPlayers(&world, *cp)

	// Draw game
	rl.BeginMode2D(gameCam.Camera2D)
	// Game scene
	rl.DrawTextureRec(gameMap.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(gameMap.Texture.Width), Height: float32(-gameMap.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
	// rl.DrawTextureRec(gameMap.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32()},)

	// Player
	// rl.DrawRectangleRec(rl.Rectangle{X: float32(cp.InternalPlayer.GetPos().Position.X) * c.BLOCK_SIZE, Y: float32(cp.InternalPlayer.GetPos().Position.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
	rl.EndMode2D()

	// // Game mouse input handler
	// gameViewRec := &rl.Rectangle{X: float32(SIDE_PANEL_WIDTH), Y: 0, Width: float32(c.SCREEN.Width) - float32(2*SIDE_PANEL_WIDTH), Height: float32(c.SCREEN.Height)}
	// mouse.HandleMouseInputs(cp, gameCam, gameViewRec, mouse.GAME)

	// // Render panels
	// rightPanel.DrawContent()
	// rightPanel.RenderPanel(*c.SCREEN)

	// leftPanel.DrawContent(cp, chatInput, chatMsgs)
	// leftPanel.RenderPanel()

	// // Test coords
	// DrawCoordinates(*cp, float32(SIDE_PANEL_WIDTH), 10)

	// // Handle inputs
	// chatInput, inputHold = keyboard.HandleChatInput(leftPanel.ActiveTab, chatInput, inputHold)
}

func DrawWorld() {}

// Read map from text file and return it as usable texture
func ImportMapFromFile(path string) error {
	// Read map file
	bytes, err := os.ReadFile("../client/assets/maps/spawn")
	if err != nil {
		fmt.Println("Can't find map file. Check your path")
	}
	// Read map size
	mapSizeEndChar := ";"
	index := strings.Index(string(bytes), mapSizeEndChar)
	if index == -1 {
		return fmt.Errorf("Couldn't read map size from file. Is the map format correct?")
	}
	// TODO: add support for asymmetrical map sizes
	mapSize, err := strconv.Atoi(string(bytes[:index]))
	if err != nil {
		return fmt.Errorf("Couldn't convert map size to integer")
	}
	fmt.Println("MAP SIZE: ", mapSize)
	layerLen := (mapSize * mapSize) * 2 // including commas
	mapContent := string(bytes[index+1:])
	layerCount := len(mapContent) / layerLen // cut in half to get real number of layers (don't count commas)
	fmt.Println("LAYERS: ", layerCount)
	// Check if all layers are present in file
	expectedLayersContentLen := layerCount * layerLen
	validLayersContent := expectedLayersContentLen == len(mapContent)
	if !validLayersContent {
		return fmt.Errorf("Expected layers content don't match file content. Layers: %d\nExpected length of layers content: %d, Got: %d", layerCount, expectedLayersContentLen, len(mapContent))
	}
	// Init texture for map
	mapSizePx := int32(int(c.BLOCK_SIZE) * mapSize)
	// mapTxt := rl.NewTexture2D(0, mapSizePx, mapSizePx, 1, rl.CompressedAstc4x4Rgba)
	// mapRt2d := rl.NewRenderTexture2D(0, mapTxt, mapTxt)
	mapTxt := rl.LoadRenderTexture(mapSizePx, mapSizePx)
	gameMap = mapTxt
	rl.BeginTextureMode(gameMap)
	for i := 0; i < layerCount-1; i++ {
		layerContent := mapContent[i*layerLen : (i+1)*layerLen]
		content := strings.Split(layerContent, ",")
		fmt.Println("LAYER CONTENT: ", content)
		for j := 0; j < len(content)-1; j++ {
			blockId, err := strconv.Atoi(content[j])
			if err != nil {
				fmt.Println("Err convert block id. Err: ", err.Error())
			}
			fmt.Println(blockId)
			block := rl.Rectangle{
				X:      float32(blockId * int(c.BLOCK_SIZE)),
				Y:      0,
				Width:  c.BLOCK_SIZE,
				Height: c.BLOCK_SIZE,
			}
			blockPos := rl.Vector2{
				X: float32((j % mapSize) * int(c.BLOCK_SIZE)),
				Y: float32((j / mapSize) * int(c.BLOCK_SIZE)),
			}
			rl.DrawTextureRec(blocksTxt, block, blockPos, rl.White)
		}
	}
	rl.EndTextureMode()
	return nil
}

// func DrawMap(world *rl.RenderTexture2D, txt *rl.Texture2D, zoom float32) {
// 	rl.BeginTextureMode(*world)
// 	rl.DrawTextureRec(*txt, rl.Rectangle{X: 0, Y: 0, Width: float32(txt.Width) * zoom, Height: float32(txt.Height) * zoom}, rl.Vector2{X: 0, Y: 0}, rl.White)
// 	rl.EndTextureMode()
// }

// Draw all online players on world texture
func DrawPlayers(world *rl.RenderTexture2D, p player.ClientPlayer) {
	cheMap := p.InternalPlayer.GetCache()
	rl.BeginTextureMode(*world)
	for _, player := range *cheMap {
		pos := player.Curr.Position
		rl.DrawRectangleRec(rl.Rectangle{X: float32(pos.X) * c.BLOCK_SIZE, Y: float32(pos.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Purple)
	}
	rl.EndTextureMode()
}

// (Tmp) draws current coordinates on the map
func DrawCoordinates(p player.ClientPlayer, x float32, y float32) {
	pPos := p.InternalPlayer.GetPos().Position
	rl.DrawTextEx(c.DEFAULT_FONT, fmt.Sprintf("<%d, %d>", pPos.X, pPos.Y), rl.Vector2{X: x, Y: y}, c.DEFAULT_FONT_SIZE, 0, rl.Blue)
}

// Draw single block on texture
// TODO: add support for Y
func DrawBlockOnTxt(blockTxt *rl.Texture2D, blockId int, x float32, y float32) {
	blockRec := rl.Rectangle{X: float32(blockId) * c.BLOCK_SIZE, Y: 0, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}
	rl.DrawTextureRec(*blockTxt, blockRec, rl.Vector2{X: x * c.BLOCK_SIZE, Y: y * c.BLOCK_SIZE}, rl.White)
}

func UnloadScene() {
	// Should it unload Texture2D made off of images too(?)
	// rl.UnloadRenderTexture(world)
}
