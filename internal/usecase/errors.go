package usecase

import "errors"

var (
	// ErrInvalidFEN is returned when the provided FEN string is malformed or invalid.
	ErrInvalidFEN = errors.New("invalid FEN string")

	// ErrSearchFailed is returned when the engine encounters a critical error during search.
	ErrSearchFailed = errors.New("search failed critically")

	// ErrNoLegalMoves is returned when the engine is asked to search a position with no legal moves.
	ErrNoLegalMoves = errors.New("no legal moves available")

	// ErrUnsupportedEngineType is returned when the requested engine type is not recognized.
	ErrUnsupportedEngineType = errors.New("unsupported engine type")

	// ErrNilSearcher is returned when a nil searcher is provided to NewGame or Game.FindBestMove.
	ErrNilSearcher = errors.New("searcher cannot be nil")
)
