package ws

// MessageType identifies the kind of WebSocket payload.
type MessageType string

// Incoming client message types.
const (
	MessageTypeInitGame      MessageType = "init_game"
	MessageTypeCalculateMove MessageType = "calculate_move"
)

// Outgoing server message types.
const (
	MessageTypeReady    MessageType = "ready"
	MessageTypeBestMove MessageType = "best_move"
	MessageTypeError    MessageType = "error"
)

// ClientMessage represents an incoming WebSocket command from a client.
type ClientMessage struct {
	Type          MessageType `json:"type"`
	FEN           string      `json:"fen,omitempty"`
	RemainingTime int         `json:"remaining_time,omitempty"` // Milliseconds.
	Increment     int         `json:"increment,omitempty"`      // Milliseconds.
	EngineType    string      `json:"engine_type,omitempty"`
	ModelID       int         `json:"model_id,omitempty"`
}

// ServerMessage represents an outgoing WebSocket event sent to a client.
type ServerMessage struct {
	Type       MessageType `json:"type"`
	BestMove   string      `json:"best_move,omitempty"`   // UCI notation (e.g., "g1f3").
	EngineType string      `json:"engine_type,omitempty"` // Only populated for MessageTypeReady.
	ModelID    int         `json:"model_id,omitempty"`    // Only populated for MessageTypeReady.
	Error      string      `json:"error,omitempty"`
}
