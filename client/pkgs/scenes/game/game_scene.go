package game

import (
	"fmt"
	c "ticoma/client/pkgs/constants"
	"ticoma/types"

	"ticoma/client/pkgs/camera"
	game_map "ticoma/client/pkgs/game_map"
	"ticoma/client/pkgs/input/keyboard"
	"ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"

	left "ticoma/client/pkgs/scenes/game/panels/left"
	right "ticoma/client/pkgs/scenes/game/panels/right"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game variables
var chatInput []byte
var chatMsgs []string
var inputHold int
var gameCam *camera.GameCamera

// Textures, assets
var blockSprite rl.Texture2D
var gameMap *game_map.GameMap

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

	// Load & Generate map from file
	gameMap = game_map.New()
	// gameMap.Init("../client/assets/maps/path_tester", &blockSprite)
	gameMap.Init("../client/assets/maps/spawn", &blockSprite)

	// Setup res, scaling
	SIDE_PANEL_WIDTH = int32((c.SCREEN.Width / 4))

	// Setup game
	gameCam = camera.New(float32(gameMap.Txt.Texture.Width/2), float32(gameMap.Txt.Texture.Height/2), float32(c.SCREEN.Width/2), float32((c.SCREEN.Height)/2))

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

	rl.BeginMode2D(gameCam.Camera2D)

	// Game map
	gameMap.Render()

	// Player
	playerRec := rl.Rectangle{
		X:      float32(cp.InternalPlayer.GetPos().Position.X-1) * c.BLOCK_SIZE,
		Y:      float32(cp.InternalPlayer.GetPos().Position.Y-1) * c.BLOCK_SIZE,
		Width:  c.BLOCK_SIZE,
		Height: c.BLOCK_SIZE,
	}

	playerNick := rl.MeasureTextEx(c.DEFAULT_FONT, *cp.Nickname, c.DEFAULT_FONT_SIZE*1.5, 0)
	rl.DrawRectangleRec(playerRec, rl.Gray)
	rl.DrawRectangleLinesEx(playerRec, 2, rl.DarkGray)
	rl.DrawTextEx(c.DEFAULT_FONT, *cp.Nickname, rl.Vector2{X: playerRec.X + (playerRec.Width / 2) - playerNick.X/2, Y: playerRec.Y - playerNick.Y}, c.DEFAULT_FONT_SIZE*1.5, 0, c.COLOR_PANEL_TEXT)
	rl.EndMode2D()

	// // Game mouse input handler
	gameViewRec := &rl.Rectangle{X: float32(SIDE_PANEL_WIDTH), Y: 0, Width: float32(c.SCREEN.Width) - float32(2*SIDE_PANEL_WIDTH), Height: float32(c.SCREEN.Height)}
	mouse.HandleMouseInputs(cp, gameMap.PathMap, gameCam, gameViewRec, mouse.GAME)

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

// Draw all online players on world texture
// func DrawPlayers(world *rl.RenderTexture2D, cp player.ClientPlayer) {
// 	cheMap := cp.InternalPlayer.GetCache()
// 	rl.BeginTextureMode(*world)
// 	for _, player := range *cheMap {
// 		pos := player.Curr.Position
// 		rl.DrawRectangleRec(rl.Rectangle{X: float32(pos.X) * c.BLOCK_SIZE, Y: float32(pos.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Purple)
// 	}
// 	rl.EndTextureMode()
// }

// (Tmp) draws current coordinates on the map
func DrawCoordinates(p player.ClientPlayer, x float32, y float32) {
	pPos := p.InternalPlayer.GetPos().Position
	rl.DrawTextEx(c.DEFAULT_FONT, fmt.Sprintf("<%d, %d>", pPos.X, pPos.Y), rl.Vector2{X: x, Y: y}, c.DEFAULT_FONT_SIZE*2, 0, c.COLOR_PANEL_OUTLINE)
}

func UnloadScene() {
	sceneReady = false
	rl.UnloadTexture(blockSprite)
	rl.UnloadRenderTexture(gameMap.Txt)
}

func HandleChatRequest(cp *player.ClientPlayer, chReq *types.ChatMessage) {
	chatterNickname := cp.InternalPlayer.GetNickname(&chReq.PeerID)
	formattedChatMsg := *chatterNickname + ": " + chReq.Message
	chatMsgs = append(chatMsgs, formattedChatMsg)
	// Clear chatInput buffer if msg was sent by Player
	if *cp.InternalPlayer.GetPeerID() == chReq.PeerID {
		chatInput = nil
	}
}

func HandleMoveRequest(cp *player.ClientPlayer, mvReq *types.PlayerPosition) {
	// Prompt move queue for other players ?
}
