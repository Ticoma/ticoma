package scenes

import (
	"fmt"
	"ticoma/client/pkgs/player"
	gameScene "ticoma/client/pkgs/scenes/game"
	mainMenuScene "ticoma/client/pkgs/scenes/main_menu"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
)

type SceneHandler struct {
	GameRunning bool
}

func New() *SceneHandler {
	return &SceneHandler{
		GameRunning: false,
	}
}

func (sh *SceneHandler) HandleScene(cp *player.ClientPlayer) {
	if !sh.GameRunning {
		mainMenuScene.UnloadScene()
		gameScene.UnloadScene()
	}

	if !cp.IsOnline {
		mainMenuScene.RenderMainMenuScene(cp)
		return
	}

	gameScene.RenderGameScene(cp)
}

// Unpack cached requests and sort them based on pfx
func (sh *SceneHandler) HandleCachedRequest(cp *player.ClientPlayer, cr types.CachedRequest) {
	switch cr.Pfx {
	case security.REGISTER_PREFIX:
		if nick, ok := cr.Req.(string); ok {
			mainMenuScene.HandleRegisterRequest(cp, nick)
		}
	case security.LOGIN_PREFIX, security.LOGOUT_PREFIX, security.DELETE_ACC_PREFIX:
		return
	case security.CHAT_PREFIX:
		if chatReq, ok := cr.Req.(types.ChatMessage); ok {
			fmt.Println(fmt.Sprintf("[CLIENT] - Received Chat request: from: %s msg: %s", chatReq.PeerID, chatReq.Message))
			gameScene.HandleChatRequest(cp, &chatReq)
		}
	case security.MOVE_PREFIX:
		if mvReq, ok := cr.Req.(types.PlayerPosition); ok {
			fmt.Println(fmt.Sprintf("[CLIENT] - Received Move request: pos: %v destPos: %v", mvReq.Position, mvReq.DestPosition))
			gameScene.HandleMoveRequest(cp, &mvReq)
		}
	default:
		fmt.Println(fmt.Sprintf("[CLIENT] - Received Cached request: {pfx : \"%s\", uncastedReq: \"%s\"}", cr.Pfx, cr.Req))
	}
}
