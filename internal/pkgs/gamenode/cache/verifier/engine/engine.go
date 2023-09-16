package engine

import (
	"fmt"
	"math"
	"ticoma/internal/debug"
	"ticoma/types"
)

// EngineVerifier
type EngineVerifier struct{}

// Engine ruleset
const (
	SECONDS_IN_MS                = 1000
	MAX_BLOCKS_TRAVEL_PER_SECOND = 4
	MAX_VEL                      = float64(MAX_BLOCKS_TRAVEL_PER_SECOND) / float64(SECONDS_IN_MS)
)

// Engine verify functions

// Checks if the player traversed too many blocks in a short amount of time
func (ev *EngineVerifier) VerifyMoveVelocity(start *types.PlayerPosition, end *types.PlayerPosition) bool {
	startPos := start.Position
	endPos := end.Position
	elapsedTime := end.Timestamp - start.Timestamp
	blocksTraveledX := math.Abs(float64(endPos.X) - float64(startPos.X))
	blocksTraveledY := math.Abs(float64(endPos.Y) - float64(startPos.Y))
	blocksTraveledTotal := blocksTraveledX + blocksTraveledY

	velocity := blocksTraveledTotal / float64(elapsedTime)

	debug.DebugLog(fmt.Sprintf("[ENGINE VER] - Move vel info: Velocity: %f, Max Vel: %f, Elapsed Time: %d\n[ENGINE VER] - First mv: { pos{%d, %d}, dest{%d, %d}}, End pos: { pos{%d, %d}, dest{%d, %d}}", velocity, MAX_VEL, elapsedTime, start.Position.X, start.Position.Y, start.DestPosition.X, start.DestPosition.Y, end.Position.X, end.Position.Y, end.DestPosition.X, end.DestPosition.Y), debug.PLAYER)

	if velocity > MAX_VEL {
		msg := fmt.Sprintf("[ERR] Engine player movement verification (velocity too high).\nVelocity: %f\nMax acceptable velocity: %f\n", velocity, MAX_VEL)
		debug.DebugLog(msg, debug.PLAYER)
		return false
	}
	return true
}

// Checks if destPos of last move req matches the pos of next req
func (ev *EngineVerifier) VerifyMoveDirection(lastDestPos *types.DestPosition, pos *types.Position) bool {
	verX, verY := lastDestPos.X == pos.X, lastDestPos.Y == pos.Y
	if !verX || !verY {
		return false
	} else {
		return true
	}
}
