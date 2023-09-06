package player

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/interfaces"
	internal_player "ticoma/internal/pkgs/player"
	"time"
)

//
// Player actions package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move request) and Updates
// the client-side player position on success
//

type ClientPlayer struct {
	InternalPlayer internal_player.Player // Internal player node
	IsActiveTile   bool                   // Helpers - hoverTile, activeTile can't be nil, and 0,0 is a valid tile
	IsHoveringTile bool
	HoverTile      *interfaces.Tile // Currently hovered game tile (empty if none)
	ActiveTile     *interfaces.Tile // Clicked tile (destination of most recent move request)
	IsMoving       bool             // Input blocker
}

func New(p internal_player.Player, initPosX int, initPosY int) *ClientPlayer {
	err := p.Move(initPosX, initPosY, initPosX, initPosY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to init player! err: ", err)
	}
	return &ClientPlayer{
		InternalPlayer: p,
		IsActiveTile:   false,
		IsHoveringTile: false,
		HoverTile:      &interfaces.Tile{},
		ActiveTile:     &interfaces.Tile{},
		IsMoving:       false,
	}
}

// Move with engine-safe delay between pkgs
func (cp *ClientPlayer) MovePlayer(cam *camera.GameCamera, destX int, destY int) {
	cp.IsMoving = true
	currPos := cp.InternalPlayer.GetPos()
	// Move request (1)
	err := cp.InternalPlayer.Move(currPos.X, currPos.Y, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (1), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	// Move arrive request (2)
	err = cp.InternalPlayer.Move(destX, destY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (2), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	cp.IsMoving = false
}

// Send a message to chat
func (cp *ClientPlayer) Chat(msg []byte) error {
	err := cp.InternalPlayer.Chat(msg)
	if err != nil {
		return fmt.Errorf("[CLIENT] - Failed to send chat msg, err: %v", err)
	}
	return nil
}
