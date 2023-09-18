package scenes

import (
	"fmt"
	gameScene "ticoma/client/pkgs/drawing/scenes/game"
	mainMenuScene "ticoma/client/pkgs/drawing/scenes/main_menu"
	"ticoma/client/pkgs/player"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
)

type SceneHandler struct {
	GameRunning bool
}

func New() *SceneHandler {
	return &SceneHandler{
		GameRunning: true,
	}
}

// Render a scene based on Player & Game state
func (sh *SceneHandler) HandleScene(cp *player.ClientPlayer) {
	// Handle unloading
	if !sh.GameRunning {
		mainMenuScene.UnloadScene()
		gameScene.UnloadScene()
	}

	// Render scene based on state
	switch cp.IsOnline {
	case false:
		mainMenuScene.RenderMainMenuScene(cp)
	case true:
		gameScene.RenderGameScene(cp)
	}
}

// Unpack cached requests and sort them based on pfx
func (sh *SceneHandler) HandleCachedRequest(cr types.CachedRequest) {
	switch cr.Pfx {
	case security.CHAT_PREFIX:
		chatReq := cr.Req.(types.ChatMessage)
		fmt.Println(fmt.Sprintf("[CLIENT] - Received Chat request: from: %s msg: %s", chatReq.PeerID, chatReq.Message))
	case security.MOVE_PREFIX:
		mvReq := cr.Req.(types.PlayerPosition)
		fmt.Println(fmt.Sprintf("[CLIENT] - Received Move request: pos: %v destPos: %v", mvReq.Position, mvReq.DestPosition))
	default:
		fmt.Println(fmt.Sprintf("[CLIENT] - Received Cached request: {pfx : \"%s\", uncastedReq: \"%s\"}", cr.Pfx, cr.Req))
	}
}
