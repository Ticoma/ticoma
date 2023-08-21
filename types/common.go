package types

// Common types needed both in internal and client
// E.g: channel types and global structs

type ChatMessage struct {
	Timestamp int64  `json:"timestamp"`
	PlayerId  int    `json:"playerId"`
	Message   string `json:"message"`
}
