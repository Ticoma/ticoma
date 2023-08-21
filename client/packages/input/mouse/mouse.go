package mouse

import (
	"ticoma/client/packages/camera"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SCROLL_DIR int

const (
	UP SCROLL_DIR = iota
	DOWN
)

// Check for all mouse inputs
func HandleMouseInputs(cam *camera.GameCamera) {
	// Scroll
	scroll := rl.GetMouseWheelMoveV().Y
	if scroll != 0 {
		if scroll > 0 {
			handleScrollZoom(UP, cam)
		} else {
			handleScrollZoom(DOWN, cam)
		}
	}
}

func handleScrollZoom(dir SCROLL_DIR, cam *camera.GameCamera) {
	if dir == UP {
		if cam.Camera2D.Zoom < 1.25 {
			cam.Camera2D.Zoom += .1
		}
	} else {
		if cam.Camera2D.Zoom > 0.75 {
			cam.Camera2D.Zoom -= .1
		}
	}
}
