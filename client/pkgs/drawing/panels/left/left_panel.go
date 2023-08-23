package left

import (
	"fmt"
	"ticoma/client/pkgs/drawing/panels"
	"ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"
	"ticoma/client/pkgs/utils"

	c "ticoma/client/pkgs/constants"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LeftPanel struct {
	*panels.SidePanel
}

func New(rt2d *rl.RenderTexture2D, width float32, height float32, renderPosX float32, renderPosY float32, bgColor *rl.Color, tabs panels.Tabs) *LeftPanel {
	return &LeftPanel{
		SidePanel: panels.New(rt2d, width, height, renderPosX, renderPosY, bgColor, tabs),
	}
}

// Renders the panel to screen
func (lp *LeftPanel) RenderPanel() {
	rl.DrawTextureRec(lp.Txt.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(lp.Txt.Texture.Width), Height: float32(-lp.Txt.Texture.Height)}, rl.Vector2{X: lp.Pos.X, Y: lp.Pos.Y}, rl.White)
}

// Draw content for the active tab (content = everything except for panel navigation tabs)
func (lp *LeftPanel) DrawContent(font *rl.Font) {

	// Always draw tabs first
	lp.DrawSkeleton()
	lp.DrawPanelTabs(font, c.DEFAULT_FONT_SIZE)

	// Content
	switch lp.Tabs[lp.ActiveTab] {
	// Chat
	case lp.Tabs[0]:
		lp.DrawPanelTitle(font, c.DEFAULT_FONT_SIZE)
	case lp.Tabs[1]:
		lp.DrawPanelTitle(font, c.DEFAULT_FONT_SIZE)
	default:
		lp.DrawPanelTitle(font, c.DEFAULT_FONT_SIZE)
	}
}

// Draws a chat block with textinput and send Button
func (lp *LeftPanel) DrawChat(panel *rl.RenderTexture2D, p internal_player.Player, yPos float32, chatInput []byte, msgs []string, font *rl.Font) {

	// Reserve some space for tabs, textinput
	chatInputHeight := float32(panel.Texture.Height / 10)

	rl.BeginTextureMode(*panel)
	rl.SetMouseCursor(rl.MouseCursorDefault)

	// Draw chat block
	rl.DrawRectangleRec(rl.Rectangle{
		X:      c.SIDE_PANEL_PADDING,
		Y:      yPos + 2*c.SIDE_PANEL_PADDING,
		Width:  float32(panel.Texture.Width) - 2*c.SIDE_PANEL_PADDING,
		Height: float32(panel.Texture.Height) - chatInputHeight - yPos - 3*c.SIDE_PANEL_PADDING},
		rl.Gray)

	// Draw chat msgs
	for i, msg := range msgs {
		rl.DrawTextEx(*font, msg, rl.Vector2{
			X: 1.5 * c.SIDE_PANEL_PADDING,
			Y: float32((i * 50)) + yPos + 2.5*c.SIDE_PANEL_PADDING,
		}, c.DEFAULT_FONT_SIZE, 0, rl.Black)
	}

	// Draw chat input
	// Block
	chatContainer := rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - 2*chatInputHeight, Width: float32(panel.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: 2*chatInputHeight - c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(chatContainer, rl.Gray)
	// Textinput ctn
	textInputRec := rl.Rectangle{X: 2 * c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width * 2 / 3), Height: chatInputHeight - 3*c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(textInputRec, rl.DarkGray)
	// Send button ctn
	sendBtnRec := rl.Rectangle{X: float32(panel.Texture.Width)*2/3 + 2*c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width*1/3) - 4*c.SIDE_PANEL_PADDING, Height: chatInputHeight - 3*c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(sendBtnRec, rl.Black)

	// Draw send button
	textSize := rl.MeasureTextEx(*font, "Send", c.DEFAULT_FONT_SIZE, 0)
	rl.DrawTextEx(*font, "Send", rl.Vector2{X: (float32(panel.Texture.Width)*2/3 + 2*c.SIDE_PANEL_PADDING) + 0.5*(float32(panel.Texture.Width*1/3)-4*c.SIDE_PANEL_PADDING) - textSize.X/2, Y: (float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING) + 0.5*(chatInputHeight-3*c.SIDE_PANEL_PADDING) - textSize.Y/2}, c.DEFAULT_FONT_SIZE, 0, rl.White)

	sendBtnHover := mouse.IsMouseHoveringRec(&sendBtnRec)

	// Handle click while focusing over chat textinput
	if len(chatInput) > 0 {
		// Draw outline indicator
		rl.DrawRectangleLinesEx(textInputRec, 2, rl.SkyBlue)

		if sendBtnHover {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				player.Chat(p, chatInput)
			}
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			player.Chat(p, chatInput)
		}
	}

	// Draw textinput
	lp.drawChatInput(panel, &textInputRec, font, chatInput)

	rl.EndTextureMode()
}

// Draws the text inside chat input box
func (lp *LeftPanel) drawChatInput(panel *rl.RenderTexture2D, inputRec *rl.Rectangle, font *rl.Font, input []byte) {

	var textOffset float32

	rl.BeginTextureMode(*panel)

	// Since the font is monospace, we can use any char here
	charSize := rl.MeasureTextEx(*font, "a", c.DEFAULT_FONT_SIZE, 0) // in px
	// Measurements
	inputRecX, centerInputRecY := inputRec.X+c.SIDE_PANEL_PADDING, inputRec.Y+0.5*inputRec.Height-0.5*charSize.Y
	inputBoxRealWidth := inputRec.Width - 3*c.SIDE_PANEL_PADDING
	maxFittingChars := int(inputBoxRealWidth / charSize.X)

	// Check if input is exceeding line width
	if len(string(input)) > maxFittingChars {
		// Draw only characters that fit in the box
		textVisible := input[len(input)-maxFittingChars-1:]
		textOffset = float32(maxFittingChars-len(textVisible)+1) * -charSize.X
		fmt.Println(string(textVisible))
		rl.DrawTextEx(*font, string(textVisible), rl.Vector2{X: inputRecX - textOffset, Y: centerInputRecY}, c.DEFAULT_FONT_SIZE, 0, rl.Black)
		// Caret
		rl.DrawTextEx(*font, "_", rl.Vector2{X: inputRecX - textOffset + float32(len(textVisible))*charSize.X, Y: centerInputRecY}, c.DEFAULT_FONT_SIZE, 0, rl.SkyBlue)
	} else {
		// Draw full textinput text
		rl.DrawTextEx(*font, string(input), rl.Vector2{X: inputRecX - textOffset, Y: centerInputRecY}, c.DEFAULT_FONT_SIZE, 0, rl.Black)
		// Caret
		rl.DrawTextEx(*font, "_", rl.Vector2{X: inputRecX - textOffset + float32(len(input))*charSize.X, Y: centerInputRecY}, c.DEFAULT_FONT_SIZE, 0, rl.SkyBlue)
	}

	rl.EndTextureMode()

}

// Temp here, just to test tabs
func DrawBuildInfo(panel *rl.RenderTexture2D, yPos float32, font *rl.Font) {

	rl.BeginTextureMode(*panel)

	hash := utils.GetCommitHash()[:6]
	rl.DrawTextEx(*font, "ticoma git-"+hash, rl.Vector2{X: c.SIDE_PANEL_PADDING, Y: yPos + c.SIDE_PANEL_PADDING}, c.DEFAULT_FONT_SIZE, 0.25, rl.Black)

	rl.EndTextureMode()
}
