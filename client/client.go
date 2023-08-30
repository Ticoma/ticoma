package client

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	c "ticoma/client/pkgs/constants"
	dr "ticoma/client/pkgs/drawing"

	"ticoma/client/pkgs/drawing/panels/left"
	"ticoma/client/pkgs/drawing/panels/right"

	"ticoma/client/pkgs/input/keyboard"
	"ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"

	"ticoma/client/pkgs/utils"
	internal_player "ticoma/internal/packages/player"
	"ticoma/types"

	"github.com/fstanis/screenresolution"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var chatInput []byte
var chatMsgs []string
var hold int

func Main(pc chan internal_player.Player, cc chan types.ChatMessage, fullscreen *bool) {

	// Load icon
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Load screen conf
	res := screenresolution.GetPrimary()
	screenC := utils.GetScreenConf(res.Width, res.Height, fullscreen)

	// Setup window, resolution, raylib
	rl.InitWindow(screenC.Width, screenC.Height, "Ticoma Client")
	if *fullscreen {
		rl.ToggleFullscreen()
	}

	defer rl.CloseWindow()
	rl.SetTraceLogLevel(rl.LogError)
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Load fonts, imgs
	c.DEFAULT_FONT = rl.LoadFontEx("../client/assets/fonts/ponderosa.regular.ttf", int32(c.DEFAULT_FONT_SIZE)*4, nil)
	spawnImg := rl.LoadImage("../client/assets/textures/map/spawn.png")
	blocksImg := rl.LoadImage("../client/assets/textures/map/blocks.png")

	// Setup res, scaling
	SIDE_PANEL_WIDTH := int32((screenC.Width / 4))

	// Setup game
	gameCam := camera.New(float32(spawnImg.Width/2), float32(spawnImg.Height/2), float32(screenC.Width/2), float32((screenC.Height+8)/2)) // where is this 8px offset coming from??????
	spawnMapSize := 25                                                                                                                    // In blocks (tmp)

	// Setup textures, panels
	world := rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height) // Full sized map

	// tmp, Draw map on world from texture
	spawnTxt := rl.LoadTextureFromImage(spawnImg)
	blocksTxt := rl.LoadTextureFromImage(blocksImg)

	// Wait for player conn
	p := <-pc

	// Activate chat listener
	go ChatMsgListener(p, cc)

	// Init ClientPlayer
	clientPlayer := player.New(p, spawnMapSize/2, spawnMapSize/2)

	// Init side panels
	leftTabs := map[int][2]string{
		0: {"Chat", "C"},
		1: {"Build info", "B"},
		2: {"Tabssss bro", "Tb"},
	}
	leftPanel := left.New(float32(SIDE_PANEL_WIDTH), float32(screenC.Height), 0, 0, &c.COLOR_PANEL_BG, leftTabs)

	rightTabs := map[int][2]string{
		0: {"Inventory", "I"},
		1: {"Settings", "S"},
	}
	rightPanel := right.New(float32(SIDE_PANEL_WIDTH), float32(screenC.Height), float32(int32(screenC.Width)-SIDE_PANEL_WIDTH), 0, &c.COLOR_PANEL_BG, rightTabs)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw players
		dr.DrawMap(&world, &spawnTxt, gameCam.Zoom)
		dr.DrawPlayers(&world, *clientPlayer)

		// Draw game
		rl.BeginMode2D(gameCam.Camera2D)
		// Game scene
		rl.DrawTextureRec(world.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(world.Texture.Width), Height: float32(-world.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		// Player
		rl.DrawRectangleRec(rl.Rectangle{X: float32(p.GetPos().X) * c.BLOCK_SIZE, Y: float32(p.GetPos().Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
		// Test block
		dr.DrawBlock(&blocksTxt, 3, 14, 14)
		rl.EndMode2D()

		// Game mouse input handler
		gameViewRec := &rl.Rectangle{X: float32(SIDE_PANEL_WIDTH), Y: 0, Width: float32(screenC.Width) - float32(2*SIDE_PANEL_WIDTH), Height: float32(screenC.Height)}
		mouse.HandleMouseInputs(clientPlayer, gameCam, gameViewRec, mouse.GAME)

		// Render panels
		rightPanel.DrawContent()
		rightPanel.RenderPanel(*screenC)

		leftPanel.DrawContent(clientPlayer, chatInput, chatMsgs)
		leftPanel.RenderPanel()

		// Test coords
		dr.DrawCoordinates(*clientPlayer, float32(SIDE_PANEL_WIDTH), 10)

		// Handle inputs
		chatInput, hold = keyboard.HandleChatInput(leftPanel.ActiveTab, chatInput, hold)

		rl.EndDrawing()
	}
}

func ChatMsgListener(p internal_player.Player, cc chan types.ChatMessage) {
	for {
		chat := <-cc
		msg := fmt.Sprintf("[player %d]: %s", chat.PlayerId, chat.Message)
		chatMsgs = append(chatMsgs, msg)
		if p.GetId() == chat.PlayerId { // if the msg arrived from us - clear input
			chatInput = nil
		}
	}
}
