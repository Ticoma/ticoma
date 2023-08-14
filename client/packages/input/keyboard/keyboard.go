package keyboard

import (
	"ticoma/client/packages/actions"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboardMoveInput(p internal_player.Player, playerMoveState *bool) {
	if *playerMoveState {
		return
	}
	pos := p.GetPos()
	x, y := pos.X, pos.Y

	if rl.IsKeyDown(rl.KeyA) {
		actions.MovePlayer(p, playerMoveState, x, y, x-1, y)
	}
	if rl.IsKeyDown(rl.KeyD) {
		actions.MovePlayer(p, playerMoveState, x, y, x+1, y)
	}
	if rl.IsKeyDown(rl.KeyS) {
		actions.MovePlayer(p, playerMoveState, x, y, x, y+1)
	}
	if rl.IsKeyDown(rl.KeyW) {
		actions.MovePlayer(p, playerMoveState, x, y, x, y-1)
	}
}
