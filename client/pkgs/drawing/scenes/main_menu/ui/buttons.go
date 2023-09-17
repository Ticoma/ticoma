package ui

import (
	c "ticoma/client/pkgs/constants"

	"ticoma/client/pkgs/input/mouse"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	Txt       *rl.RenderTexture2D
	Pos       rl.Vector2
	BgColor   *rl.Color // Background
	BdColor   *rl.Color // Border
	Text      string
	OnClick   func()
	ColOffset rl.Vector2
}

func NewButton(rt2d *rl.RenderTexture2D, width float32, height float32, x float32, y float32, bgColor *rl.Color, bdColor *rl.Color, text string, onClick func()) *Button {
	return &Button{
		Txt: rt2d,
		Pos: rl.Vector2{
			X: x,
			Y: y,
		},
		BgColor:   bgColor,
		BdColor:   bdColor,
		Text:      text,
		OnClick:   onClick,
		ColOffset: rl.Vector2{},
	}
}

func (btn *Button) Draw() {
	rl.BeginTextureMode(*btn.Txt)
	btnTextSize := rl.MeasureTextEx(c.DEFAULT_FONT, btn.Text, c.DEFAULT_FONT_SIZE*2, 0)
	btnRec := rl.Rectangle{X: 0, Y: 0, Width: float32(btn.Txt.Texture.Width), Height: float32(btn.Txt.Texture.Height)}
	rl.DrawRectangleRec(btnRec, *btn.BgColor)
	rl.DrawRectangleLinesEx(btnRec, 2, *btn.BdColor)
	rl.DrawTextEx(c.DEFAULT_FONT, btn.Text, rl.Vector2{X: btnRec.X + btnRec.Width/2 - btnTextSize.X/2, Y: btnRec.Y + btnRec.Height/2 - btnTextSize.Y/2}, c.DEFAULT_FONT_SIZE*2, 0, c.COLOR_PANEL_TEXT)
	rl.EndTextureMode()
}

func (btn *Button) Render() {
	rl.DrawTextureRec(btn.Txt.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(btn.Txt.Texture.Width), Height: -float32(btn.Txt.Texture.Height)}, rl.Vector2{X: btn.Pos.X, Y: btn.Pos.Y}, rl.White)
	btn.ColOffset.X, btn.ColOffset.Y = btn.Pos.X, btn.Pos.Y // Update collision offset
}

func (btn *Button) HandleClick(inputEnabled bool) bool {
	btnColRec := rl.Rectangle{
		X:      btn.ColOffset.X,
		Y:      btn.ColOffset.Y,
		Width:  float32(btn.Txt.Texture.Width),
		Height: float32(btn.Txt.Texture.Height),
	}
	if inputEnabled {
		if mouse.IsMouseHoveringRec(&btnColRec) {
			rl.DrawRectangleLinesEx(btnColRec, 2, c.COLOR_PANEL_OUTLINE)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				btn.Text = "..."
				btn.OnClick()
				return false
			}
		}
	}
	return inputEnabled
}
