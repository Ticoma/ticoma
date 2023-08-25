package interfaces

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScreenInfo struct {
	MonitorId   int
	Width       int32
	Height      int32
	RefreshRate int
}

type Tile struct {
	rl.Vector2
}
