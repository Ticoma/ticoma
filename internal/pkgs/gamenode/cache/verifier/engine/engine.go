package engine

import (
	"fmt"
	"math"
	"strconv"
	"ticoma/internal/debug"
	intf "ticoma/internal/pkgs/gamenode/cache/interfaces"
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
func (ev *EngineVerifier) VerifyPlayerMovement(startPkg *intf.ActionDataPackageTimestamped, endPkg *intf.ActionDataPackageTimestamped) bool {

	startPos := startPkg.Position
	endPos := endPkg.Position
	elapsedTime := endPkg.Timestamp - startPkg.Timestamp
	blocksTraveledX := math.Abs(float64(endPos.X) - float64(startPos.X))
	blocksTraveledY := math.Abs(float64(endPos.Y) - float64(startPos.Y))
	blocksTraveledTotal := blocksTraveledX + blocksTraveledY

	velocity := blocksTraveledTotal / float64(elapsedTime)

	debug.DebugLog("[ENGINE] VELOCITY: "+strconv.FormatFloat(velocity, 'E', -1, 32), debug.PLAYER)
	debug.DebugLog("[ENGINE] MAX VELOCITY: "+strconv.FormatFloat(MAX_VEL, 'E', -1, 32), debug.PLAYER)
	debug.DebugLog("[ENGINE] ELAPSED TIME: "+strconv.FormatInt(elapsedTime, 10), debug.PLAYER)

	if velocity > MAX_VEL {
		msg := fmt.Sprintf("[ERR] Engine player movement verification (velocity too high)\nVelocity: %f, max acceptable velocity: %f", velocity, MAX_VEL)
		debug.DebugLog(msg, debug.PLAYER)
		return false
	}

	return true

}

// Checks if destPos of last package matches the pos of currently arriving package
func (ev *EngineVerifier) VerifyLastMovePos(lastDestPos *intf.DestPosition, pos *intf.Position) bool {
	verX := lastDestPos.X == pos.X
	verY := lastDestPos.Y == pos.Y
	if !verX || !verY {
		return false
	} else {
		return true
	}
}
