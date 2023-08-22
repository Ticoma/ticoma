package left

import (
	"ticoma/client/packages/player"

	c "ticoma/client/packages/constants"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draws a chat block with textinput and send Button
//
// Returns [textInputActive state, textInputRect]
func DrawChat(panel *rl.RenderTexture2D, p internal_player.Player, yPos float32, chatInput []byte, msgs []string, font *rl.Font) (bool, *rl.Rectangle) {

	chatInputHeight := float32(panel.Texture.Height / 10) // Reserve some space for input

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

	// Hovering over textinput (blue)
	if textInputActive {
		rl.DrawRectangleLinesEx(textInputRec, 2, rl.SkyBlue)
		rl.SetMouseCursor(rl.MouseCursorIBeam)
		// Enter press while typing
		if rl.IsKeyPressed(rl.KeyEnter) {
			player.Chat(p, chatInput)
		}
	}

	// Hovering over send button (green)
	if sendBtnActive {
		rl.DrawRectangleLinesEx(sendBtn, 2, rl.Green)
		// On click send button
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			player.Chat(p, chatInput)
		}
	}

	rl.EndTextureMode()
	return textInputActive, &textInputRec
}

// Draws the text inside chat input box
func DrawChatInputText(panel *rl.RenderTexture2D, inputRec *rl.Rectangle, font *rl.Font, input []byte) {

	rl.BeginTextureMode(*panel)

	textSize := rl.MeasureTextEx(*font, string(input), c.DEFAULT_FONT_SIZE, 0)
	inputRecX, centerInputRecY := inputRec.X+c.SIDE_PANEL_PADDING, inputRec.Y+0.5*inputRec.Height-0.5*textSize.Y

	rl.DrawTextEx(*font, string(input), rl.Vector2{X: inputRecX, Y: centerInputRecY}, c.DEFAULT_FONT_SIZE, 0, rl.Black)

	rl.EndTextureMode()

}
