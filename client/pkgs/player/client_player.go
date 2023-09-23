package player

import (
	"fmt"
	c "ticoma/client/pkgs/constants"
	internal_player "ticoma/internal/pkgs/player"
	"ticoma/types"
	"time"

	"github.com/solarlune/paths"
)

//
// Player actions package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move request) and Updates
// the client-side player position on success
//

type ClientPlayer struct {
	IsOnline        bool
	InternalPlayer  internal_player.Player // Interface passed from internal
	HoverTile       *types.Position        // Currently hovered game tile
	DestTile        *types.Position        // Destination of most recent move request
	MoveQueue       []paths.Cell           // Queued DestPositions to reach DestTile
	MoveQueueCancel chan bool              // Channel for stopping an ongoing MouseMove
}

func New(p *internal_player.Player) *ClientPlayer {
	return &ClientPlayer{
		IsOnline:        false,
		InternalPlayer:  *p,
		HoverTile:       nil,
		DestTile:        nil,
		MoveQueue:       []paths.Cell{},
		MoveQueueCancel: make(chan bool),
	}
}

// Send a REGISTER_ request from Client
func (cp *ClientPlayer) Register() {
	nick := "teaver"
	err := cp.InternalPlayer.Register(&nick)
	if err != nil {
		fmt.Println("[CLIENT PLAYER] - Failed to register. Err: " + err.Error())
		return
	}
	cp.IsOnline = true
}

// Send a LOGIN_ request from Client
func (cp *ClientPlayer) Login() {
	err := cp.InternalPlayer.Login()
	if err != nil {
		fmt.Println("[CLIENT PLAYER] - Failed to login. Err: " + err.Error())
		return
	}
	cp.IsOnline = true
}

// Computes a path and queues n MOVE_ requests to destination tile
func (cp *ClientPlayer) MouseMove(cMap *paths.Grid) error {

	if cp.DestTile == nil {
		return fmt.Errorf("[CLIENT PLAYER] - MouseMove err: No Destination is set")
	}

	if len(cp.MoveQueue) != 0 {
		cp.MoveQueueCancel <- true
	}

	path := cMap.GetPath(float64(cp.InternalPlayer.GetPos().Position.X-1), float64(cp.InternalPlayer.GetPos().Position.Y-1), float64(cp.DestTile.X-1), float64(cp.DestTile.Y-1), false, false)

	if path != nil {

		ticker := time.NewTicker(time.Millisecond * time.Duration(c.MOVE_COOLDOWN_IN_MS))

		// Fill queue with new path tiles
		for i, pathTile := range path.Cells {
			// Paths pkg counts first pos, but we don't need it. MoveQueue stores only Dest's
			if i != 0 {
				cp.MoveQueue = append(cp.MoveQueue, *pathTile)
			}
		}

		// fmt.Println("Queue: ", cp.MoveQueue)

		// Perform move in
		for range cp.MoveQueue {
			go func() {
				select {
				case <-cp.MoveQueueCancel:
					// New MouseMove request - Abort this one and clear queue
					cp.MoveQueue = nil
					return
				case <-ticker.C:

					if len(cp.MoveQueue) <= 0 {
						return
					}

					// fmt.Println("Move: ", cp.MoveQueue[0].X+1, cp.MoveQueue[0].Y+1)

					queuedX, queuedY := cp.MoveQueue[0].X+1, cp.MoveQueue[0].Y+1

					// First move (Intention)
					mvErr := cp.InternalPlayer.Move(&cp.InternalPlayer.GetPos().Position.X, &cp.InternalPlayer.GetPos().Position.Y, &queuedX, &queuedY)
					if mvErr != nil {
						fmt.Println("[CLIENT PLAYER] - Failed to MouseMove (Intention) Err: ", mvErr.Error())
						return
					}

					// Second move (Arrival)
					mvErr = cp.InternalPlayer.Move(&queuedX, &queuedY, &queuedX, &queuedY)
					if mvErr != nil {
						fmt.Println("[CLIENT PLAYER] - Failed to MouseMove (Arrival) Err: ", mvErr.Error())
						return
					}

					// This panics when I change move direction in certain tick timings
					// TODO: Solution: Make a global tickrate and make a Wait on cancel

					cp.MoveQueue = cp.MoveQueue[1:] // Pop first element from queue
				}
			}()
		}
	}
	return nil
}
