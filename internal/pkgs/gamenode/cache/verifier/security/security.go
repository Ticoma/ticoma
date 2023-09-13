package security

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/cache/utils"
	"ticoma/types"
	"time"
)

// SecurityVerifier
type SecurityVerifier struct{}

// Request schemas and prefixes
const (

	//
	// Account-related
	//

	REGISTER_PREFIX   = "REGISTER_"
	DELETE_ACC_PREFIX = "DELETEACC_"
	LOGIN_PREFIX      = "LOGIN_"
	LOGOUT_PREFIX     = "LOGOUT_"

	//
	// Game-related
	//

	MOVE_PREFIX = "MOVE_"
	MOVE_SCHEMA = `{
		pos: {
			posX: int,
			posY: int,
		},
		destPos: {
			destPosX: int,
			destPosY: int,
		},
	}`
	CHAT_PREFIX = "CHAT_"
	CHAT_SCHEMA = `{
		peerId: string,
		message: string,
	}`
)

// Checks if there's a prefix in req data and if so, returns it
func (sv *SecurityVerifier) DetectReqPrefix(reqData []byte) (string, error) {
	prefixChar := "_"
	index := strings.Index(string(reqData), prefixChar)
	if index == -1 {
		return "", fmt.Errorf("[SEC VER] - Prefix not found in request.%s", "")
	}
	prefix := string(reqData[:index+1])
	return prefix, nil
}

// Returns a json schema of a request
func (sv *SecurityVerifier) getReqSchema(reqPrefix string) (string, error) {
	switch reqPrefix {
	case MOVE_PREFIX:
		return MOVE_SCHEMA, nil
	case CHAT_PREFIX:
		return CHAT_SCHEMA, nil
	case REGISTER_PREFIX, DELETE_ACC_PREFIX, LOGIN_PREFIX, LOGOUT_PREFIX:
		return "", nil
	default:
		return "", fmt.Errorf("[SEC VER] - Couldn't find schema of request with this prefix: %s", reqPrefix)
	}
}

// Verifies package types and returns Req data as string (no prefix)
func (sv *SecurityVerifier) VerifyReqTypes(prefix string, reqData []byte) (string, error) {

	schema, err := sv.getReqSchema(prefix)
	if err != nil {
		return "", fmt.Errorf("[SEC VER] - Err while getting req schema: %v", err)
	}

	reqStr := strings.TrimPrefix(string(reqData), prefix)
	res := []byte{}
	keySelected := false

	dec := json.NewDecoder(strings.NewReader(reqStr))

	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				// fmt.Println(err)
				debug.DebugLog(fmt.Sprintf("Err while decoding req: %v", err), debug.PLAYER)
			}
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

	debug.DebugLog("[SEC VER] - Schema types: "+utils.StripString(schema, false), debug.PLAYER)
	debug.DebugLog("[SEC VER] - Req types: "+string(res), debug.PLAYER)

	valid := strings.Compare(utils.StripString(schema, false), utils.StripString(string(res), true)) == 0
	if !valid {
		return "", fmt.Errorf("[SEC VER] - Couldn't validate request data with schema. %s", "")
	}
	return reqStr, nil
}

// Request interface from raw pubsub data
func (sv *SecurityVerifier) ReqFromBytes(peerID *string, data *[]byte) (types.Request, error) {

	if peerID == nil || data == nil {
		return types.Request{}, fmt.Errorf("[SEC VER] - Null peerID or reqData pointer")
	}

	if len(*peerID) == 0 || len(*data) == 0 {
		return types.Request{}, fmt.Errorf("[SEC VER] - PeerID or reqData is empty?")
	}

	req := types.Request{
		PeerID: *peerID,
		Data:   *data,
	}
	return req, nil
}

// construct request based on prefix
func (sv *SecurityVerifier) AutoConstructRequest(prefix string, data string) (interface{}, error) {
	switch prefix {
	case MOVE_PREFIX:
		moveReq, err := sv.constructMoveReq(data)
		return moveReq, err
	case CHAT_PREFIX:
		chatReq, err := sv.constructChatReq(data)
		return chatReq, err
	default:
		return nil, nil
	}
}

// construct PlayerPosition from req data
func (sv *SecurityVerifier) constructMoveReq(data string) (types.PlayerPosition, error) {

	const EXPECTED_VAL_LEN_IN_MOVE = 4 // [posX, posY, destPosX, destPosY]
	var IGNORED_STRINGS_IN_MOVE = []string{"posX", "posY", "pos", "destPosX", "destPosY", "destPos"}

	vals := utils.ExtractValsFromStrPkg(data, IGNORED_STRINGS_IN_MOVE)
	if len(vals) != EXPECTED_VAL_LEN_IN_MOVE {
		return types.PlayerPosition{Position: &types.Position{}, DestPosition: &types.DestPosition{}}, fmt.Errorf("[SEC VER] - Couldn't extract - pkg values length don't match schema")
	}

	var positions []int
	for i := 0; i < len(vals); i++ {
		pos, err := strconv.Atoi(vals[i])
		if err != nil {
			return types.PlayerPosition{Position: &types.Position{}, DestPosition: &types.DestPosition{}}, fmt.Errorf("[SEC VER] - Err while converting string val to int")
		}
		positions = append(positions, pos)
	}

	ts := time.Now().UnixMilli()
	playerPos := types.PlayerPosition{
		Timestamp:    ts,
		Position:     &types.Position{X: positions[0], Y: positions[1]},
		DestPosition: &types.DestPosition{X: positions[2], Y: positions[3]},
	}
	return playerPos, nil
}

// construct ChatMessage from req data
func (sv *SecurityVerifier) constructChatReq(data string) (types.ChatMessage, error) {

	const EXPECTED_VAL_LEN_IN_CHAT = 2 // [message, peerId]
	var IGNORED_STRINGS_IN_CHAT = []string{"message", "peerId"}

	vals := utils.ExtractValsFromStrPkg(data, IGNORED_STRINGS_IN_CHAT)
	if len(vals) != EXPECTED_VAL_LEN_IN_CHAT {
		return types.ChatMessage{}, fmt.Errorf("[SEC VER] - Couldn't extract - pkg values length don't match schema")
	}

	ts := time.Now().UnixMilli()
	chatPkg := types.ChatMessage{
		PeerID:    vals[0],
		Timestamp: ts,
		Message:   vals[1],
	}
	return chatPkg, nil
}
