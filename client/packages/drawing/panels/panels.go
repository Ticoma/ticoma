package panels

import (
	c "ticoma/client/packages/constants"
	"ticoma/client/packages/player"

	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//
// Package responsible for drawing side panels and any
// items that belong inside them
//

// Draws a panel skeleton (solid color) at pos {x, y}
func DrawSidePanelSkeleton(panel *rl.RenderTexture2D, width float32, height float32, col *rl.Color) {
	rl.BeginTextureMode(*panel)

	rl.DrawRectangleRec(rl.Rectangle{X: 0, Y: 0, Width: width, Height: height}, *col)

	rl.EndTextureMode()
}

// Draw title at specified Y pos on panel and returns its height
func DrawTitleBlock(panel *rl.RenderTexture2D, yPos float32, title string, font *rl.Font) float32 {

	titleBlockHeight := float32(panel.Texture.Height / 12)
	titleSize := rl.MeasureTextEx(*font, title, c.DEFAULT_FONT_SIZE, 0)

	rl.BeginTextureMode(*panel)

	rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: yPos + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width - 2*c.SIDE_PANEL_PADDING), Height: titleBlockHeight - float32(2*c.SIDE_PANEL_PADDING)}, rl.Gray)
	rl.DrawTextEx(*font, title, rl.Vector2{X: float32(panel.Texture.Width/2) - titleSize.X/2, Y: titleBlockHeight/2 - titleSize.Y/2}, c.DEFAULT_FONT_SIZE, 0, rl.Black)

	rl.EndTextureMode()

	return titleBlockHeight - float32(2*c.SIDE_PANEL_PADDING)
}

// Draws a chat block with textinput and send Button
//
// Returns [newChatInput, newMsgsArr, textInputActive state, textInputRect]
func DrawChat(panel *rl.RenderTexture2D, p internal_player.Player, yPos float32, chatInput []byte, msgs []string, font *rl.Font) (bool, *rl.Rectangle) {

	chatInputHeight := float32(panel.Texture.Height / 10)

	rl.BeginTextureMode(*panel)

	rl.SetMouseCursor(rl.MouseCursorDefault)

	// Draw chat block
	rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: yPos + 2*c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: float32(panel.Texture.Height) - chatInputHeight - yPos - 3*c.SIDE_PANEL_PADDING}, rl.Gray)

	// Draw chat msgs
	for i, msg := range msgs {
		rl.DrawTextEx(*font, msg, rl.Vector2{
			X: 1.5 * c.SIDE_PANEL_PADDING,
			Y: float32((i * 50)) + yPos + 2.5*c.SIDE_PANEL_PADDING,
		}, c.DEFAULT_FONT_SIZE, 0, rl.Black)
	}

	// Draw chat input
	// Block
	chatContainer := rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - chatInputHeight, Width: float32(panel.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: chatInputHeight - c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(chatContainer, rl.Gray)
	// Textinput ctn
	textInputRec := rl.Rectangle{X: 2 * c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width * 2 / 3), Height: chatInputHeight - 3*c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(textInputRec, rl.DarkGray)
	// Send button ctn
	sendBtn := rl.Rectangle{X: float32(panel.Texture.Width)*2/3 + 3*c.SIDE_PANEL_PADDING, Y: float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width*1/3) - 5*c.SIDE_PANEL_PADDING, Height: chatInputHeight - 3*c.SIDE_PANEL_PADDING}
	rl.DrawRectangleRec(sendBtn, rl.DarkGray)

	// Draw send button
	textSize := rl.MeasureTextEx(*font, "Send", c.DEFAULT_FONT_SIZE, 0)
	rl.DrawTextEx(*font, "Send", rl.Vector2{X: (float32(panel.Texture.Width)*2/3 + 3*c.SIDE_PANEL_PADDING) + 0.5*(float32(panel.Texture.Width*1/3)-5*c.SIDE_PANEL_PADDING) - textSize.X/2, Y: (float32(panel.Texture.Height) - chatInputHeight + c.SIDE_PANEL_PADDING) + 0.5*(chatInputHeight-3*c.SIDE_PANEL_PADDING) - textSize.Y/2}, c.DEFAULT_FONT_SIZE, 0, rl.White)

	// Return chat active state (is user hovering textinput)
	textInputActive := rl.CheckCollisionPointRec(rl.GetMousePosition(), textInputRec)
	sendBtnActive := rl.CheckCollisionPointRec(rl.GetMousePosition(), sendBtn)

	// draw outlines if active, react to send inputs
	if textInputActive {
		rl.DrawRectangleLinesEx(textInputRec, 2, rl.SkyBlue)
		rl.SetMouseCursor(rl.MouseCursorIBeam)
		// Enter press while typing
		if rl.IsKeyPressed(rl.KeyEnter) {
			player.Chat(p, chatInput)
		}
	}
	if sendBtnActive {
		rl.DrawRectangleLinesEx(sendBtn, 2, rl.Green)
		// On click send button
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			player.Chat(p, chatInput)
		}
	}

	rl.EndTextureMode()

	// return states
	return textInputActive, &textInputRec
}

func DrawChatInput(panel *rl.RenderTexture2D, textPos *rl.Rectangle, font *rl.Font, input []byte) {

	rl.BeginTextureMode(*panel)

	textPosX, centerTextPosY := textPos.X+c.SIDE_PANEL_PADDING, textPos.Y+0.5*textPos.Height
	rl.DrawTextEx(*font, string(input), rl.Vector2{X: textPosX, Y: centerTextPosY}, c.DEFAULT_FONT_SIZE, 0, rl.Black)

	rl.EndTextureMode()

}
