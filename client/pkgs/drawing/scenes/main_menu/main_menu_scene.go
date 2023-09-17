package mainmenu

import (
	c "ticoma/client/pkgs/constants"
	"ticoma/client/pkgs/drawing/scenes/main_menu/ui"
	"ticoma/client/pkgs/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var regBtn *ui.Button
var logBtn *ui.Button
var inputEnabled = true
var sceneReady bool = false

func RenderMainMenuScene(cp *player.ClientPlayer) {
	if !sceneReady {
		loadMainMenuScene(cp)
	}

	welcomeTextSize := rl.MeasureTextEx(c.DEFAULT_FONT, "Ticoma", c.DEFAULT_FONT_SIZE*6, 0)
	welcomeTextPos := rl.Vector2{
		X: float32(c.SCREEN.Width/2) - (welcomeTextSize.X / 2),
		Y: float32(c.SCREEN.Height/5) - (welcomeTextSize.Y / 2),
	}
	rl.DrawTextEx(c.DEFAULT_FONT, "Ticoma", rl.Vector2{X: welcomeTextPos.X, Y: welcomeTextPos.Y}, c.DEFAULT_FONT_SIZE*6, 0, rl.White)

	regBtn.Draw()
	logBtn.Draw()
	regBtn.Render()
	logBtn.Render()
	inputEnabled = regBtn.HandleClick(inputEnabled)
	inputEnabled = logBtn.HandleClick(inputEnabled)
}

func loadMainMenuScene(cp *player.ClientPlayer) {

	var btnWidth int32 = 350
	var btnHeight int32 = 100
	var btnGap int32 = 30 // gap between a button and scene center
	regBtnPos := rl.Vector2{
		X: float32((c.SCREEN.Width / 2) - btnGap - btnWidth),
		Y: float32(c.SCREEN.Height / 2),
	}
	logBtnPos := rl.Vector2{
		X: float32((c.SCREEN.Width / 2) + btnGap),
		Y: float32(c.SCREEN.Height / 2),
	}
	regBtnTxt := rl.LoadRenderTexture(btnWidth, btnHeight)
	logBtnTxt := rl.LoadRenderTexture(btnWidth, btnHeight)
	regBtn = ui.NewButton(&regBtnTxt, float32(regBtnTxt.Texture.Width), float32(regBtnTxt.Texture.Height), regBtnPos.X, regBtnPos.Y, &c.COLOR_PANEL_BG, &c.COLOR_PANEL_CONTENT, "Register", cp.Register)
	logBtn = ui.NewButton(&logBtnTxt, float32(regBtnTxt.Texture.Width), float32(regBtnTxt.Texture.Height), logBtnPos.X, logBtnPos.Y, &c.COLOR_PANEL_BG, &c.COLOR_PANEL_CONTENT, "Login (SOON)", cp.Login)

	sceneReady = true
}

func UnloadScene() {
	rl.UnloadRenderTexture(*regBtn.Txt)
	rl.UnloadRenderTexture(*logBtn.Txt)
}
