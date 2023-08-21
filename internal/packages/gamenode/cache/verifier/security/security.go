package security

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"ticoma/internal/debug"
	"ticoma/internal/packages/gamenode/cache/interfaces"
	"ticoma/internal/packages/gamenode/cache/utils"
	"ticoma/types"
	"time"
)

// SecurityVerifier
type SecurityVerifier struct{}

type PACKAGE_TYPE int

const (
	ADP PACKAGE_TYPE = iota
	CHAT
)

// Returns a json schema of a package
func (sv *SecurityVerifier) getSchemaAndPrefix(reqSchema PACKAGE_TYPE) (string, string, error) {

	// Package schema collection
	// const ADPPrefix = "ADP_" // TODO
	const schemaADP = `{
		playerId: int,
		pubKey: string,
		pos: {
			posX: int,
			posY: int,
		},
		destPos: {
			destPosX: int,
			destPosY: int,
		},
	},`

	const chatPrefix = "CHAT_"
	const schemaChat = `{
		playerId: int,
		message: string,
	},`

	switch reqSchema {
	case ADP:
		return schemaADP, "", nil
	case CHAT:
		return schemaChat, chatPrefix, nil
	default:
		return "", "", fmt.Errorf("[SEC VER] - Couldn't find requested schema")
	}
}

// Verifies package types and returns
// Str pkg without a prefix
func (sv *SecurityVerifier) verifyPackageTypes(pkg []byte, pkgType PACKAGE_TYPE) (bool, string) {

	schema, prefix, err := sv.getSchemaAndPrefix(pkgType)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	pkgStr := string(pkg)

	// Trim prefix if needed
	pkgStr = strings.TrimPrefix(pkgStr, prefix)

	res := []byte{}
	keySelected := false

	// Anti spam
	if len(pkg) == 0 {
		return false, ""
	}

	dec := json.NewDecoder(strings.NewReader(pkgStr))

	for {
		t, err := dec.Token()
		if err != nil {
			break
		}

		switch v := t.(type) {
		case json.Delim:
			if string(v) == "}" || string(v) == "]" {
				res = append(res, []byte(string(v)+", ")...)
			} else {
				res = append(res, byte(v))
			}
			if keySelected {
				keySelected = false
			}
		case string:
			if !keySelected {
				keySelected = true
				res = append(res, []byte(v+": ")...)
			} else {
				res = append(res, []byte(fmt.Sprintf("%T, ", v))...)
				keySelected = false
			}
		case float64:
			res = append(res, []byte(fmt.Sprintf("%T, ", int(v)))...)
			keySelected = false
		default:
			res = append(res, []byte(fmt.Sprintf("%T, ", v))...)
			keySelected = false
		}
	}

	debug.DebugLog("[ADP TYPES] SCHEMA "+utils.StripString(schema, true), debug.PLAYER)
	debug.DebugLog("[ADP TYPES RES "+utils.StripString(string(res), true), debug.PLAYER)

	valid := strings.Compare(utils.StripString(schema, true), utils.StripString(string(res), true)) == 0

	return valid, pkgStr
}

// Verify pkg types and construct ADPT struct from byte pkg
func (sv *SecurityVerifier) ConstructADPT(pkgBytes []byte) (interfaces.ActionDataPackageTimestamped, error) {

	const EXPECTED_VAL_LENGTH_IN_ADP = 6 // [playerId, pubKey, posX, posY, destX, destY]

	// Type check
	validPkgTypes, _ := sv.verifyPackageTypes(pkgBytes, ADP)
	if !validPkgTypes {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't verify package types")
	}

	// If types are OK, try extract vals
	vals := utils.ExtractValsFromStrPkg(string(pkgBytes))
	if len(vals) != EXPECTED_VAL_LENGTH_IN_ADP {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't extract - pkg values length don't match schema")
	}

	fmt.Println(vals[0])
	playerId, err := strconv.Atoi(vals[0])
	if err != nil {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't assign playerId from vals")
	}

	var positions []int
	// conv vals[2:5] to ints
	for i := 2; i < len(vals); i++ {
		pos, err := strconv.Atoi(vals[i])
		if err != nil {
			return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Err while converting string val to int")
		}
		positions = append(positions, pos)
	}

	timestamp := time.Now().UnixMilli()

	ADPT := interfaces.ActionDataPackageTimestamped{
		ActionDataPackage: &interfaces.ActionDataPackage{
			PlayerId:     playerId,
			PubKey:       vals[1],
			Position:     &interfaces.Position{X: positions[0], Y: positions[1]},
			DestPosition: &interfaces.DestPosition{X: positions[2], Y: positions[3]},
		},
		Timestamp: timestamp,
	}

	return ADPT, nil
}

// Verify pkg types and construct ChatPkg from byte pkg
func (sv *SecurityVerifier) ConstructChatPkg(pkgBytes []byte) (types.ChatMessage, error) {

	const EXPECTED_VAL_LENGTH_IN_CHAT_PKG = 2 // [playerId, message]

	validPkgTypes, pkgStr := sv.verifyPackageTypes(pkgBytes, CHAT)
	if !validPkgTypes {
		return types.ChatMessage{}, fmt.Errorf("[SEC VER] - Couldn't verify package types")
	}

	// If types are OK, try extract vals
	vals := utils.ExtractValsFromStrPkg(pkgStr)
	if len(vals) != EXPECTED_VAL_LENGTH_IN_CHAT_PKG {
		return types.ChatMessage{}, fmt.Errorf("[SEC VER] - Couldn't extract - pkg values length don't match schema")
	}

	playerId, err := strconv.Atoi(vals[0])
	if err != nil {
		return types.ChatMessage{}, fmt.Errorf("[SEC VER] - Couldn't assign playerId from pkg vals")
	}

	timestamp := time.Now().UnixMilli()

	chatPkg := types.ChatMessage{
		Timestamp: timestamp,
		PlayerId:  playerId,
		Message:   vals[1],
	}

	debug.DebugLog(fmt.Sprintf("[SEC VER] - Chat pkg constructed! pkg: %+v\n", chatPkg), debug.PLAYER)

	return chatPkg, nil
}
