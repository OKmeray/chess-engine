package ws

import "errors"

// Errors for the WebSocket transport layer.
var (
	ErrCreatorNotConfigured = errors.New("game creator is not configured")
	ErrInvalidMessage       = errors.New("invalid message format")
	ErrEngineBusy           = errors.New("engine is busy calculating another move")
	ErrInitGameFailed       = errors.New("failed to initialize game")
	ErrEngineNotInitialized = errors.New("engine is not initialized")
	ErrCalculationFailed    = errors.New("failed to calculate move")
)
