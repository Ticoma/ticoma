package types

//
// Common types needed both in internal and client
//

// Generic network request
type Request struct {
	PeerID string
	Data   []byte
}

// Result of Cache's Put and verified by Cache's components.
// Sent directly to Client from Gamenode
type CachedRequest struct {
	Pfx string
	Req interface{}
}

type Position struct {
	X int `json:"posX"`
	Y int `json:"posY"`
}

type DestPosition struct {
	X int `json:"destPosX"`
	Y int `json:"destPosY"`
}

type Player struct {
	IsOnline       bool `json:"isOnline"`
	PlayerGameData `json:"playerGameData"`
}

type PlayerGameData struct {
	Nick string
	PlayerPosition
}

type PlayerPosition struct {
	Timestamp    int64 `json:"timestamp"`
	Position     `json:"pos"`
	DestPosition `json:"destPos"`
}

type ChatMessage struct {
	PeerID    string `json:"peerId"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}
