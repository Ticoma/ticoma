package keyboard

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// func HandleKeyboardMoveInput(p internal_player.Player, cam *camera.GameCamera, playerMoveState *bool) {
// 	if *playerMoveState {
// 		return
// 	}
// 	pos := p.GetPos()
// 	x, y := pos.X, pos.Y

// 	if rl.IsKeyDown(rl.KeyA) {
// 		player.MovePlayer(p, cam, camera.LEFT, playerMoveState, x, y, x-1, y)
// 		cam.Target.X -= c.BLOCK_SIZE
// 	}
// 	if rl.IsKeyDown(rl.KeyD) {
// 		player.MovePlayer(p, cam, camera.RIGHT, playerMoveState, x, y, x+1, y)
// 		cam.Target.X += c.BLOCK_SIZE
// 	}
// 	if rl.IsKeyDown(rl.KeyS) {
// 		player.MovePlayer(p, cam, camera.DOWN, playerMoveState, x, y, x, y+1)
// 		cam.Target.Y += c.BLOCK_SIZE
// 	}
// 	if rl.IsKeyDown(rl.KeyW) {
// 		player.MovePlayer(p, cam, camera.UP, playerMoveState, x, y, x, y-1)
// 		cam.Target.Y -= c.BLOCK_SIZE
// 	}
// }

// Handle key presses when chat is active (appends to value of provided input byte arr)
// Increments hold value if backspace is held down
func HandleChatInput(activeTabLeft int, input []byte, hold int) ([]byte, int) {

	key := rl.GetCharPressed()

	// Handle chat behavior
	if activeTabLeft == 0 { // Chat is open TODO: This is bad. Make a global enum called inputMode with chat, prompt modes, etc.
		// Ignore space when empty input
		if len(input) == 0 && rl.IsKeyPressed(rl.KeySpace) {
			return input, hold
		}

		// TODO: add more character support
		if key != 0 {
			if key >= 32 && key <= 125 {
				input = append(input, byte(key))
			}
		}

		// Handle backspace behavior
		// Single press
		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		}
		// Hold
		if rl.IsKeyDown(rl.KeyBackspace) {
			hold++
			if len(input) > 0 && hold >= 10 { // 10 frames hold = fast delete
				input = input[:len(input)-1]
				time.Sleep(time.Millisecond * 40) // without 40ms cd it is a bit too fast
			}
		} else {
			hold = 0
		}
	}

	return input, hold
}
