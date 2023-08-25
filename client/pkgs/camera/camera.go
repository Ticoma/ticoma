package camera

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SCROLL_DIR int

const (
	UP SCROLL_DIR = iota
	DOWN
)

type GameCamera struct {
	rl.Camera2D
}

func New(targetX float32, targetY float32, offsetX float32, offsetY float32) *GameCamera {
	offset := rl.Vector2{X: offsetX, Y: offsetY}
	target := rl.Vector2{X: targetX, Y: targetY}
	return &GameCamera{
		Camera2D: rl.NewCamera2D(offset, target, 0, 1),
	}
}

// Zoom camera in/out on mouse wheel
func HandleScrollZoom(dir SCROLL_DIR, cam *GameCamera) {
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
