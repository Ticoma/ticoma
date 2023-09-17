package game

import (
	"fmt"
	c "ticoma/client/pkgs/constants"

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
var world rl.RenderTexture2D
var spawnImg *rl.Image
var spawnTxt rl.Texture2D
var blocksTxt rl.Texture2D
var blocksImg *rl.Image

// Side panels
var leftPanel *left.LeftPanel
var rightPanel *right.RightPanel

// Misc
var sceneReady bool = false
var SIDE_PANEL_WIDTH int32

// Load all the textures, assets needed to render game scene
func loadGameScene() {
	spawnImg = rl.LoadImage("../client/assets/textures/map/spawn.png")
	blocksImg = rl.LoadImage("../client/assets/textures/map/blocks.png")

	// Setup res, scaling
	SIDE_PANEL_WIDTH = int32((c.SCREEN.Width / 4))

	// Setup game
	gameCam = camera.New(float32(spawnImg.Width/2), float32(spawnImg.Height/2), float32(c.SCREEN.Width/2), float32((c.SCREEN.Height)/2))
	// spawnMapSize := 25

	// Setup textures, panels
	world = rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height) // Full sized map

	// tmp, Draw map on world from texture
	spawnTxt = rl.LoadTextureFromImage(spawnImg)
	blocksTxt = rl.LoadTextureFromImage(blocksImg)

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

	// Draw players
	DrawMap(&world, &spawnTxt, gameCam.Zoom)
	DrawPlayers(&world, *cp)

	// Draw game
	rl.BeginMode2D(gameCam.Camera2D)
	// Game scene
	rl.DrawTextureRec(world.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(world.Texture.Width), Height: float32(-world.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
	// Player
	rl.DrawRectangleRec(rl.Rectangle{X: float32(cp.InternalPlayer.GetPos().Position.X) * c.BLOCK_SIZE, Y: float32(cp.InternalPlayer.GetPos().Position.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
	// Test block (question mark)
	DrawBlock(&blocksTxt, 3, 14, 14)
	rl.EndMode2D()

	// Game mouse input handler
	gameViewRec := &rl.Rectangle{X: float32(SIDE_PANEL_WIDTH), Y: 0, Width: float32(c.SCREEN.Width) - float32(2*SIDE_PANEL_WIDTH), Height: float32(c.SCREEN.Height)}
	mouse.HandleMouseInputs(cp, gameCam, gameViewRec, mouse.GAME)

	// Render panels
	rightPanel.DrawContent()
	rightPanel.RenderPanel(*c.SCREEN)

	leftPanel.DrawContent(cp, chatInput, chatMsgs)
	leftPanel.RenderPanel()

	// Test coords
	DrawCoordinates(*cp, float32(SIDE_PANEL_WIDTH), 10)

	// Handle inputs
	chatInput, inputHold = keyboard.HandleChatInput(leftPanel.ActiveTab, chatInput, inputHold)
}

// Draw empty world map from texture
//
// TODO: Draw map from json-like object imported from texture editor instead of img
func DrawMap(world *rl.RenderTexture2D, txt *rl.Texture2D, zoom float32) {
	rl.BeginTextureMode(*world)
	rl.DrawTextureRec(*txt, rl.Rectangle{X: 0, Y: 0, Width: float32(txt.Width) * zoom, Height: float32(txt.Height) * zoom}, rl.Vector2{X: 0, Y: 0}, rl.White)
	rl.EndTextureMode()
}

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

// (Tmp) draws a block from block sprite
func DrawBlock(blockTxt *rl.Texture2D, id int, mapX float32, mapY float32) {
	blockRec := rl.Rectangle{X: float32(id) * c.BLOCK_SIZE, Y: 0, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}
	rl.DrawTextureRec(*blockTxt, blockRec, rl.Vector2{X: mapX * c.BLOCK_SIZE, Y: mapY * c.BLOCK_SIZE}, rl.White)
}

func UnloadScene() {
	// Should it unload Texture2D made off of images too(?)
	rl.UnloadRenderTexture(world)
	rl.UnloadImage(spawnImg)
	rl.UnloadImage(blocksImg)
}
