package left

import (
	"ticoma/client/pkgs/drawing/scenes/game/panels"

	"ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"
	"ticoma/client/pkgs/utils"

	c "ticoma/client/pkgs/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LeftPanel struct {
	*panels.SidePanel
}

func New(width float32, height float32, renderPosX float32, renderPosY float32, bgColor *rl.Color, tabs panels.Tabs) *LeftPanel {
	rt2d := rl.LoadRenderTexture(int32(width), int32(height)) // Create render texture for panel
	sp := panels.New(&rt2d, width, height, renderPosX, renderPosY, bgColor, tabs)
	return &LeftPanel{
		SidePanel: &sp,
	}
}

// Renders the panel to screen
func (lp *LeftPanel) RenderPanel() {
	rl.DrawTextureRec(lp.Txt.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(lp.Txt.Texture.Width), Height: float32(-lp.Txt.Texture.Height)}, rl.Vector2{X: lp.Pos.X, Y: lp.Pos.Y}, c.COLOR_PANEL_TEXT)
}

// Draw content for the active tab (content = everything except for panel navigation tabs)
func (lp *LeftPanel) DrawContent(cp *player.ClientPlayer, chatInput []byte, chatMsgs []string) {

	// Always draw tabs first
	lp.SidePanel.DrawSkeleton()
	lp.SidePanel.DrawPanelTabs()

	// Content
	switch lp.Tabs[lp.ActiveTab] {
	// Chat
	case lp.Tabs[0]:
		// Drawing order: Tabs -> Chat input -> Title -> Chat
		lp.DrawPanelTitle()
		lp.DrawChatInput(cp, chatInput)
		lp.DrawChat(chatInput, chatMsgs)
	case lp.Tabs[1]:
		lp.DrawPanelTitle()
		lp.DrawBuildInfo()
	default:
		lp.DrawPanelTitle()
	}
}

// Draws a chat block with textinput and send Button
func (lp *LeftPanel) DrawChat(chatInput []byte, msgs []string) {

	// Measurements
	chatCtnY := lp.SidePanel.SpaceTaken.Top + c.SIDE_PANEL_PADDING
	chatCtnW := float32(lp.SidePanel.Txt.Texture.Width) - 3*c.SIDE_PANEL_PADDING
	chatCtnH := float32(lp.SidePanel.Txt.Texture.Height) - chatCtnY - lp.SidePanel.SpaceTaken.Bottom
	charSize := rl.MeasureTextEx(c.DEFAULT_FONT, "a", c.DEFAULT_FONT_SIZE, 0)
	maxCharsInLine := int(chatCtnW / charSize.X)

	rl.BeginTextureMode(*lp.SidePanel.Txt)

	// Draw chat block
	rl.DrawRectangleRec(rl.Rectangle{
		X:      c.SIDE_PANEL_PADDING,
		Y:      chatCtnY,
		Width:  chatCtnW + c.SIDE_PANEL_PADDING,
		Height: chatCtnH,
	}, c.COLOR_PANEL_CONTENT)

	// Draw chat msgs
	lines := 0
	for _, msg := range msgs {
		// Calc lines of msg
		msgLines := 1 + len(msg)/maxCharsInLine
		// Multiline
		if msgLines > 1 {
			// If long msg, chunk it
			for j := 0; j < msgLines; j++ {
				// Check how much content is left
				if len(msg[j*maxCharsInLine:]) >= maxCharsInLine {
					chunk := utils.FirstN(msg[j*maxCharsInLine:], maxCharsInLine)
					drawChatMsg(chunk, 1.5*c.SIDE_PANEL_PADDING, float32((lines*50))+lp.SidePanel.SpaceTaken.Top+2.5*c.SIDE_PANEL_PADDING)
					lines++
				} else {
					chunk := msg[j*maxCharsInLine:]
					drawChatMsg(chunk, 1.5*c.SIDE_PANEL_PADDING, float32((lines*50))+lp.SidePanel.SpaceTaken.Top+2.5*c.SIDE_PANEL_PADDING)
					lines++
				}
			}
		} else {
			// Single line msg
			drawChatMsg(msg, 1.5*c.SIDE_PANEL_PADDING, float32((lines*50))+lp.SidePanel.SpaceTaken.Top+2.5*c.SIDE_PANEL_PADDING)
			lines++
		}
	}

	rl.EndTextureMode()
}

// Draws the text inside chat input box
func (lp *LeftPanel) DrawChatInput(cp *player.ClientPlayer, chatInput []byte) {

	// Dimensions
	prevTakenB := lp.SidePanel.SpaceTaken.Bottom          // Need to update this at the end
	chatInputCtnH := lp.SidePanel.Txt.Texture.Height / 12 // (no padding)
	chatInputCtnY := float32(lp.SidePanel.Txt.Texture.Height) - prevTakenB - float32(chatInputCtnH) + c.SIDE_PANEL_PADDING
	chatInputCtnW := float32(lp.SidePanel.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING
	var textOffset float32

	rl.BeginTextureMode(*lp.SidePanel.Txt)

	// Draw chat input container (textinput + btn)
	chatInputRec := &rl.Rectangle{
		X:      c.SIDE_PANEL_PADDING,
		Y:      chatInputCtnY,
		Width:  chatInputCtnW,
		Height: float32(chatInputCtnH) - c.SIDE_PANEL_PADDING,
	}
	rl.DrawRectangleRec(*chatInputRec, c.COLOR_PANEL_CONTENT)

	// Textinput ctn
	inputRec := &rl.Rectangle{
		X:      c.SIDE_PANEL_PADDING,
		Y:      chatInputCtnY,
		Width:  chatInputCtnW * 4 / 5,
		Height: float32(chatInputCtnH) - c.SIDE_PANEL_PADDING,
	}

	// Send btn
	sendRec := &rl.Rectangle{
		X:      inputRec.Width + c.SIDE_PANEL_PADDING,
		Y:      chatInputCtnY,
		Width:  chatInputCtnW * 1 / 5,
		Height: float32(chatInputCtnH) - c.SIDE_PANEL_PADDING,
	}
	rl.DrawRectangleRec(*sendRec, rl.Black) // change this later
	// Send btn text
	sendBtnMsg := "Send"
	sendTextSize := rl.MeasureTextEx(c.DEFAULT_FONT, sendBtnMsg, c.DEFAULT_FONT_SIZE, 0)
	rl.DrawTextEx(c.DEFAULT_FONT, sendBtnMsg, rl.Vector2{
		X: sendRec.X + sendRec.Width/2 - sendTextSize.X/2,
		Y: sendRec.Y + sendRec.Height/2 - sendTextSize.Y/2,
	}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)

	// Draw border around chat if input is not empty
	if len(chatInput) > 0 {
		rl.DrawRectangleLinesEx(*inputRec, 2, c.COLOR_PANEL_OUTLINE)
	}

	// Measurements
	// Since the font is monospace, we can use any char here
	charSize := rl.MeasureTextEx(c.DEFAULT_FONT, "a", c.DEFAULT_FONT_SIZE, 0) // in px
	inputRecX, centerInputRecY := inputRec.X+c.SIDE_PANEL_PADDING, inputRec.Y+0.5*inputRec.Height-0.5*charSize.Y
	inputBoxRealWidth := inputRec.Width - 3*c.SIDE_PANEL_PADDING
	maxFittingChars := int(inputBoxRealWidth / charSize.X)

	// Textinput draw
	caret := "_"
	// Check if input is exceeding line width
	if len(string(chatInput)) > maxFittingChars {
		// Draw only characters that fit in the textinput bounds
		textVisible := chatInput[len(chatInput)-maxFittingChars-1:]
		textOffset = float32(maxFittingChars-len(textVisible)+1) * -charSize.X
		drawChatInput(string(textVisible), inputRecX-textOffset, centerInputRecY, caret, inputRecX-textOffset+float32(len(textVisible))*charSize.X, centerInputRecY)
	} else {
		// Draw full input
		drawChatInput(string(chatInput), inputRecX-textOffset, centerInputRecY, caret, inputRecX-textOffset+float32(len(chatInput))*charSize.X, centerInputRecY)
	}

	rl.EndTextureMode()

	// Handle send button click
	sendHover := mouse.IsMouseHoveringRec(sendRec)

	if sendHover {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			cp.Chat(&chatInput)
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		cp.Chat(&chatInput)
	}

	// Update space taken
	lp.SidePanel.SpaceTaken.Bottom = prevTakenB + float32(chatInputCtnH)
}

// Temp here, just to test tabs
func (lp *LeftPanel) DrawBuildInfo() {

	rl.BeginTextureMode(*lp.SidePanel.Txt)

	hash := utils.GetCommitHash()[:6]
	rl.DrawTextEx(c.DEFAULT_FONT, "ticoma client", rl.Vector2{X: c.SIDE_PANEL_PADDING, Y: lp.SidePanel.SpaceTaken.Top + c.SIDE_PANEL_PADDING}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)
	rl.DrawTextEx(c.DEFAULT_FONT, "git-"+hash, rl.Vector2{X: c.SIDE_PANEL_PADDING, Y: lp.SidePanel.SpaceTaken.Top + 40 + c.SIDE_PANEL_PADDING}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)

	rl.EndTextureMode()
}

// Draws message content at given position
func drawChatMsg(content string, xPos float32, yPos float32) {
	rl.DrawTextEx(c.DEFAULT_FONT, content, rl.Vector2{
		X: xPos,
		Y: yPos,
	}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)
}

// Draws text in chat input box (input text + caret)
func drawChatInput(content string, ctntX float32, ctntY float32, caret string, caretX float32, caretY float32) {
	// Text
	rl.DrawTextEx(c.DEFAULT_FONT, content, rl.Vector2{X: ctntX, Y: ctntY}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)
	// Caret
	rl.DrawTextEx(c.DEFAULT_FONT, caret, rl.Vector2{X: caretX, Y: caretY}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_OUTLINE)
}
