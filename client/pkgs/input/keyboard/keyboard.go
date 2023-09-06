package keyboard

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Handle key presses when chat is active (appends to value of provided input byte arr)
// Increments hold value if backspace is held down
func HandleChatInput(activeTabLeft int, input []byte, hold int) ([]byte, int) {

	key := rl.GetCharPressed()

	// Handle chat behavior
	if activeTabLeft == 0 {
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
