package game

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	c "ticoma/client/pkgs/constants"
	"ticoma/types"

	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/input/keyboard"
	"ticoma/client/pkgs/input/mouse"
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
var blockSprite rl.Texture2D
var gameMap rl.RenderTexture2D

// Side panels
var leftPanel *left.LeftPanel
var rightPanel *right.RightPanel

// Misc
var sceneReady bool = false
var SIDE_PANEL_WIDTH int32

// Load all the textures, assets needed to render game scene
func loadGameScene() {

	// Load sprites
	blockSprite = rl.LoadTexture("../client/assets/textures/sprites/blocks.png")

	err := ImportMapFromFile("../client/assets/maps/spawn")
	if err != nil {
		fmt.Println("Failed to import game map from file. Err: ", err.Error())
		return
	}

	// Setup res, scaling
	SIDE_PANEL_WIDTH = int32((c.SCREEN.Width / 4))

	// Setup game
	gameCam = camera.New(float32(gameMap.Texture.Width/2), float32(gameMap.Texture.Height/2), float32(c.SCREEN.Width/2), float32((c.SCREEN.Height)/2))

	// Init side panels
	leftTabs := map[int][2]string{
		0: {"Chat", "C"},
		1: {"Build info", "B"},
		2: {"Tabssss bro", "Tb"},
	}
	leftPanel = left.New(float32(SIDE_PANEL_WIDTH), float32(c.SCREEN.Height), 0, 0, &c.COLOR_PANEL_BG, leftTabs)

	rightTabs := map[int][2]string{
		0: {"Inventory", "I"},
		1: {"Settings", "S"},
	}
	rightPanel = right.New(float32(SIDE_PANEL_WIDTH), float32(c.SCREEN.Height), float32(int32(c.SCREEN.Width)-SIDE_PANEL_WIDTH), 0, &c.COLOR_PANEL_BG, rightTabs)

	sceneReady = true
}

// Render game scene (requires Player to be logged in)
func RenderGameScene(cp *player.ClientPlayer) {

	if !sceneReady {
		loadGameScene()
	}

	// // Draw players
	// DrawPlayers(&world, *cp)

	rl.BeginMode2D(gameCam.Camera2D)
	// Game game map
	rl.DrawTextureRec(gameMap.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(gameMap.Texture.Width), Height: float32(-gameMap.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)

	// Player
	rl.DrawRectangleRec(rl.Rectangle{X: float32(cp.InternalPlayer.GetPos().Position.X-1) * c.BLOCK_SIZE, Y: float32(cp.InternalPlayer.GetPos().Position.Y-1) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
	rl.EndMode2D()

	// // Game mouse input handler
	gameViewRec := &rl.Rectangle{X: float32(SIDE_PANEL_WIDTH), Y: 0, Width: float32(c.SCREEN.Width) - float32(2*SIDE_PANEL_WIDTH), Height: float32(c.SCREEN.Height)}
	mouse.HandleMouseInputs(cp, gameCam, gameViewRec, mouse.GAME)

	// // Render panels
	rightPanel.DrawContent()
	rightPanel.RenderPanel()

	leftPanel.DrawContent(cp, chatInput, chatMsgs)
	leftPanel.RenderPanel()

	// // Test coords
	DrawCoordinates(*cp, float32(SIDE_PANEL_WIDTH), 0)

	// // Handle inputs
	chatInput, inputHold = keyboard.HandleChatInput(leftPanel.ActiveTab, chatInput, inputHold)
}

// Load game map from file and render it on specified texture
//
// TODO: Divide this bulky func to smaller chunks (loading & interpreting file, drawing?)
func ImportMapFromFile(path string) error {
	// Read map file
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Can't find map file. Check your path")
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
	// fmt.Println("MAP SIZE: ", mapSize)
	layerLen := (mapSize * mapSize) * 2 // including commas
	mapContent := string(bytes[index+1:])
	layerCount := len(mapContent) / layerLen // cut in half to get real number of layers (don't count commas)
	// Check if all layers are present in file
	expectedLayersContentLen := layerCount * layerLen
	validLayersContent := expectedLayersContentLen == len(mapContent)
	if !validLayersContent {
		return fmt.Errorf("Expected layers content don't match file content. Layers: %d\nExpected length of layers content: %d, Got: %d", layerCount, expectedLayersContentLen, len(mapContent))
	}
	// Init texture for map
	mapSizePx := int32(int(c.BLOCK_SIZE) * mapSize)
	gameMap = rl.LoadRenderTexture(mapSizePx, mapSizePx)

	rl.BeginTextureMode(gameMap)
	// Draw blocks from sprite to corresponding pos on map texture
	for i := 0; i < layerCount; i++ {
		layerContent := mapContent[i*layerLen : (i+1)*layerLen]
		content := strings.Split(layerContent, ",")
		for j := 0; j < len(content)-1; j++ {
			blockId, err := strconv.Atoi(content[j])
			if err != nil {
				return fmt.Errorf("Err convert block id. Err: %s", err.Error())
			}
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
			// Ignore block 0 (transparent)
			if blockId != 0 {
				rl.DrawTextureRec(blockSprite, block, blockPos, rl.White)
			}
		}
	}
	rl.EndTextureMode()
	return nil
}

// Draw all online players on world texture
func DrawPlayers(world *rl.RenderTexture2D, cp player.ClientPlayer) {
	cheMap := cp.InternalPlayer.GetCache()
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
	rl.DrawTextEx(c.DEFAULT_FONT, fmt.Sprintf("<%d, %d>", pPos.X, pPos.Y), rl.Vector2{X: x, Y: y}, c.DEFAULT_FONT_SIZE*2, 0, c.COLOR_PANEL_OUTLINE)
}

func UnloadScene() {
	sceneReady = false
	rl.UnloadTexture(blockSprite)
	rl.UnloadRenderTexture(gameMap)
}

//
// Game request handler
// TODO: Think about client side request handling and info flow
//

func HandleMoveRequest() {
	// Todo
}

func HandleChatRequest(cp *player.ClientPlayer, chReq *types.ChatMessage) {
	chatMsgs = append(chatMsgs, chReq.Message)
	// Clear chatInput buffer if msg came from us
	if cp.InternalPlayer.GetPeerID() == chReq.PeerID {
		chatInput = nil
	}
}
