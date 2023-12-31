package client

import (
	"fmt"
	"ticoma/client/packages/camera"
	c "ticoma/client/packages/constants"
	dr "ticoma/client/packages/drawing"
	"ticoma/client/packages/drawing/panels"
	"ticoma/client/packages/input/keyboard"
	"ticoma/client/packages/input/mouse"
	"ticoma/client/packages/player"
	utils "ticoma/client/packages/utils"
	internal_player "ticoma/internal/packages/player"
	"ticoma/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var chatMsgs []string
var chatInput []byte

func Main(pc chan internal_player.Player, cc chan types.ChatMessage) {

	// Load icon
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Setup window, raylib
	rl.InitWindow(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), "Ticoma Client")
	// rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Load fonts, imgs
	font := rl.LoadFont("../client/assets/fonts/ponderosa.regular.ttf")
	spawnImg := rl.LoadImage("../client/assets/textures/map/spawn.png")

	// Setup res, scaling
	screenConf := utils.GetScreenInfo()
	SIDE_PANEL_WIDTH := int32((screenConf.Width / 4))

	// Setup game
	gameCam := camera.New(float32(spawnImg.Width/2), float32(spawnImg.Height/2), float32(screenConf.Width/2), float32(screenConf.Height/2))
	playerMoved := false
	spawnMapSize := 25 // In blocks

	// Setup textures, panels
	world := rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height)                 // Full sized map
	rightPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height)) // Side panels
	leftPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))

	// tmp, Draw map on world from texture
	spawnTxt := rl.LoadTextureFromImage(spawnImg)

	// Wait for player conn
	p := <-pc

	// Activate chat listener
	go ChatMsgListener(cc)

	// Init player
	player.InitPlayer(p, &playerMoved, spawnMapSize/2, spawnMapSize/2)

	// Draw textures that only need to be drawn once
	panelColor := rl.DarkGray
	panels.DrawSidePanelSkeleton(&leftPanel, float32(SIDE_PANEL_WIDTH), float32(screenConf.Height), &panelColor)
	panels.DrawSidePanelSkeleton(&rightPanel, float32(SIDE_PANEL_WIDTH), float32(screenConf.Height), &panelColor)
	chatTitleH := panels.DrawTitleBlock(&leftPanel, 0, "Chat", &font)
	panels.DrawTitleBlock(&rightPanel, 0, "7357", &font)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw players
		dr.DrawMap(&world, &spawnTxt, gameCam.Zoom)
		dr.DrawPlayers(&world, p, gameCam.Zoom)

		// Draw chat
		chatActive, textInputRec := panels.DrawChat(&leftPanel, p, chatTitleH, chatInput, chatMsgs, &font)
		panels.DrawChatInput(&leftPanel, textInputRec, &font, chatInput)

		// Draw game
		rl.BeginMode2D(gameCam.Camera2D)
		rl.DrawTextureRec(world.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(world.Texture.Width), Height: float32(-world.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawRectangleRec(rl.Rectangle{X: float32(p.GetPos().X * c.BLOCK_SIZE), Y: float32(p.GetPos().Y * c.BLOCK_SIZE), Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
		rl.EndMode2D()

		// Draw panels
		rl.DrawTextureRec(leftPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(leftPanel.Texture.Width), Height: float32(-leftPanel.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawTextureRec(rightPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(rightPanel.Texture.Width), Height: float32(-rightPanel.Texture.Height)}, rl.Vector2{X: float32(int32(screenConf.Width) - SIDE_PANEL_WIDTH), Y: 0}, rl.White)

		// Ignore game keyboard inputs if user is typing in chat
		if !chatActive {
			keyboard.HandleKeyboardMoveInput(p, gameCam, &playerMoved)
		} else {
			chatInput = keyboard.HandleChatInput(chatInput)
		}
		mouse.HandleMouseInputs(gameCam)

		rl.EndDrawing()
	}

}

func ChatMsgListener(cc chan types.ChatMessage) {
	for {
		chat := <-cc
		msg := fmt.Sprintf("[ player %d ]: %s", chat.PlayerId, chat.Message)
		chatMsgs = append(chatMsgs, msg)
		chatInput = nil
	}
}
