package player

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/interfaces"
	internal_player "ticoma/internal/packages/player"
	"time"
)

//
// Player package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move) and Updates
// the client-side player position
//

type ClientPlayer struct {
	InternalPlayer internal_player.Player // Internal player node
	IsHoveringTile bool                   // Helpers - hoverTile, activeTile can't be nil, and 0,0 is a valid tile
	IsActiveTile   bool
	HoverTile      *interfaces.Tile // Currently hovered game tile (empty if )
	ActiveTile     *interfaces.Tile // Clicked tile (destination of interaction)
	IsMoving       bool             // Input blocker
}

func New(p internal_player.Player, initPosX int, initPosY int) *ClientPlayer {
	err := p.Move(initPosX, initPosY, initPosX, initPosY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to init player! err: ", err)
	}
	return &ClientPlayer{
		InternalPlayer: p,
		IsHoveringTile: false,
		IsActiveTile:   false,
		HoverTile:      &interfaces.Tile{},
		ActiveTile:     &interfaces.Tile{},
		IsMoving:       false,
	}
}

// Like move, but single pkg -> it should fill both places in empty cache
// func InitPlayer(p internal_player.Player, moveState *bool, posX int, posY int) {
// 	*moveState = true
// 	err := p.Move(posX, posY, posX, posY)
// 	if err != nil {
// 		fmt.Println("[CLIENT] - Failed to init player! err: ", err)
// 	}
// 	time.Sleep(time.Millisecond * 300)
// 	*moveState = false
// }

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
	// Move arrive (2)
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
