package scenes

import (
	gameScene "ticoma/client/pkgs/drawing/scenes/game"
	mainMenuScene "ticoma/client/pkgs/drawing/scenes/main_menu"
	"ticoma/client/pkgs/player"
)

// tmp name, can't think of a better one now
type ProgramState struct {
	Running bool
}

func New() *ProgramState {
	return &ProgramState{
		Running: true,
	}
}

func (gs *ProgramState) HandleScene(cp *player.ClientPlayer) {
	// Handle unloading
	if !gs.Running {
		mainMenuScene.UnloadScene()
	}

	// Render scene based on state
	switch cp.IsOnline {
	case false:
		mainMenuScene.RenderMainMenuScene(cp)
	case true:
		gameScene.RenderGameScene(cp)
	}
}
