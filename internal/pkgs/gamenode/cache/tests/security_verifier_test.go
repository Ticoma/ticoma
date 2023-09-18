package tests

import (
	"fmt"
	"testing"
	"ticoma/internal/pkgs/gamenode/cache/verifier"
	"ticoma/types"

	assert "github.com/stretchr/testify/assert"
)

var testPeerID string = "1230022446600"
var nv verifier.NodeVerifier

func TestReqConstruct(t *testing.T) {

	nv = *verifier.New()

	//
	// Test basic manual request construction from byte data
	//

	randomBytes := []byte("just_a_random_string")
	req, err := nv.SecurityVerifier.ReqFromBytes(&testPeerID, &randomBytes)
	if err != nil {
		t.Fatal("Failed create req from bytes: ", err)
	}
	expected := types.Request{
		PeerID: testPeerID,
		Data:   randomBytes,
	}

	// Even though the data is nonsense, construct itself should pass
	assert.Equal(t, expected.Data, req.Data)
	assert.Equal(t, expected.PeerID, req.PeerID)

	// Pass empty peerID
	_, err2 := nv.SecurityVerifier.ReqFromBytes(nil, &randomBytes)

	// Should throw pointer to nil
	assert.Error(t, err2)

	emptyStr := ""
	_, err3 := nv.SecurityVerifier.ReqFromBytes(&emptyStr, &randomBytes)

	// Should throw empty req
	assert.Error(t, err3)

}

func TestReqPrefixDetection(t *testing.T) {

	//
	// Cache should be able to detect prefixes in requests
	//

	testMsg := "Hello there"
	unknownReqData := []byte(`TEST_{invalid:"jsondata"}`)
	noPrefixReqData := []byte(`_{no:"prefix"}`)
	randomReqData := []byte(`++#$#z{awdawdwad}''....\`)

	validMoveReqData := []byte(fmt.Sprintf(`MOVE_{pos: {posX:%d,posY:%d,},destPos: {destPosX:%d,destPosY:%d,},},`, 1, 1, 2, 2))
	validChatReqData := []byte(fmt.Sprintf(`CHAT_{message:"%s",},`, testMsg))

	_, err0 := nv.DetectReqPrefix(unknownReqData)
	_, err1 := nv.DetectReqPrefix(noPrefixReqData)
	_, err2 := nv.DetectReqPrefix(validMoveReqData)
	_, err3 := nv.DetectReqPrefix(validChatReqData)
	_, errRand := nv.DetectReqPrefix(randomReqData)

	// Should pass
	assert.NoError(t, err0)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	// Expect the no-prefix request to fail
	expectedErr0 := "[SEC VER] - Prefix not found in request."
	assert.EqualErrorf(t, errRand, expectedErr0, "Error should be: %v, got: %v", expectedErr0, errRand)
}

func TestReqTypesVerification(t *testing.T) {

	//
	// Cache should detect wrong value types in requests
	//

	msg := "Hello from tests"
	moveReqPrefix := "MOVE_"
	moveReqData := []byte(fmt.Sprintf(`{"pos": {"posX":%d,"posY":%d},"destPos": {"destPosX":%d,"destPosY":%d}}`, 1, 1, 2, 2))
	chatReqPrefix := "CHAT_"
	chatReqData := []byte(fmt.Sprintf(`{"message":"%s"}`, msg))

	_, err := nv.VerifyReqTypes(moveReqPrefix, moveReqData)
	_, err2 := nv.VerifyReqTypes(chatReqPrefix, chatReqData)

	assert.ErrorIs(t, err, nil)
	assert.ErrorIs(t, err2, nil)

	// Invalid reqs (wrong types)

	invalidMoveReqData := []byte(fmt.Sprintf(`{"pos": {"posX":%d,"posY":%d},"destPos": {"destPosX":%d,"destPosY":%v}}`, 1, 1, 2, nil))
	invalidChatReqData := []byte(fmt.Sprintf(`{"message":%t}`, true))

	_, err3 := nv.VerifyReqTypes(moveReqPrefix, invalidMoveReqData)
	_, err4 := nv.VerifyReqTypes(chatReqPrefix, invalidChatReqData)
	assert.Error(t, err3)
	assert.Error(t, err4)

}

func TestAutoReqConstruct(t *testing.T) {

	//
	// Verifier should be able to construct a full request based on data
	// Using all the modules tested above
	//

	moveReqPrefix := "MOVE_"
	moveReqData := fmt.Sprintf(`{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 1, 1, 2, 2)
	expectedMoveReq := types.PlayerPosition{
		Timestamp:    0,
		Position:     types.Position{X: 1, Y: 1},
		DestPosition: types.DestPosition{X: 2, Y: 2},
	}
	moveReqIntf, err := nv.AutoConstructRequest(moveReqPrefix, moveReqData, playerID)
	moveReq := moveReqIntf.(types.PlayerPosition)
	// "Delete" the timestamp since we can't recreate it
	moveReq.Timestamp = 0

	chatReqPrefix := "CHAT_"
	chatReqData := fmt.Sprintf(`{"message":%s}`, "SomeMsg")
	chatReqIntf, err2 := nv.AutoConstructRequest(chatReqPrefix, chatReqData, playerID)
	chatReq := chatReqIntf.(types.ChatMessage)
	// "Delete" the timestamp since we can't recreate it
	chatReq.Timestamp = 0

	expectedChatReq := types.ChatMessage{
		PeerID:    playerID,
		Timestamp: 0,
		Message:   "SomeMsg",
	}

	// Should pass
	assert.NoError(t, err)
	assert.Equal(t, expectedMoveReq, moveReq)
	assert.NoError(t, err2)
	assert.Equal(t, expectedChatReq, chatReq)

}
