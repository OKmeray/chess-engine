package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

const (
	defaultGameLength = 40
	minRemainingMoves = 20                     // Capping divisor in late endgame
	safetyMargin      = 500 * time.Millisecond // Safety buffer to prevent losing on time
	minSearchTime     = 100 * time.Millisecond
)

// MoveFinder finds the best move for a chess position within time constraints.
type MoveFinder interface {
	FindBestMove(ctx context.Context, fen string, remainingTime, increment time.Duration) (string, error)
}

// Game orchestrates the business logic for the chess engine.
type Game struct {
	searcher Searcher
}

// NewGame creates a new Game with the provided searcher.
func NewGame(searcher Searcher) (*Game, error) {
	if searcher == nil {
		return nil, ErrNilSearcher
	}
	return &Game{searcher: searcher}, nil
}

// FindBestMove analyzes a given FEN string and returns the calculated UCI move.
func (g *Game) FindBestMove(
	ctx context.Context,
	fen string,
	remainingTime, increment time.Duration,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if g == nil || g.searcher == nil {
		return "", ErrNilSearcher
	}

	pos, err := engine.ParseFEN(fen)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidFEN, err)
	}

	var movesBuf [256]engine.Move
	legalMoves := pos.GenerateMoves(movesBuf[:0])
	if len(legalMoves) == 0 {
		return "", ErrNoLegalMoves
	}

	if len(legalMoves) == 1 {
		return legalMoves[0].UCI(), nil
	}

	targetTime := calculateTargetTime(pos.CurrentTurn, len(legalMoves), remainingTime, increment)
	searchCtx, cancel := context.WithTimeout(ctx, targetTime)
	defer cancel()

	bestMove, err := g.searcher.Search(searchCtx, pos)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrSearchFailed, err)
	}

	if bestMove == 0 {
		return "", fmt.Errorf("%w: searcher returned empty move", ErrSearchFailed)
	}

	if !slices.Contains(legalMoves, bestMove) {
		return "", fmt.Errorf("%w: searcher returned illegal move %s", ErrSearchFailed, bestMove.UCI())
	}

	return bestMove.UCI(), nil
}

func calculateTargetTime(currentTurn, numLegalMoves int, remainingTime, increment time.Duration) time.Duration {
	turns := max(currentTurn, 1)
	expectedRemaining := max(defaultGameLength-turns, minRemainingMoves)

	base := (remainingTime / time.Duration(expectedRemaining)) + (increment * 3 / 4)

	// Complexity scaling based on position branching factor
	if numLegalMoves > 45 {
		base = (base * 5) / 4 // 1.25x for complex positions
	} else if numLegalMoves < 15 {
		base = (base * 3) / 4 // 0.75x for simple positions
	}

	target := max(base, minSearchTime)

	hardLimit := remainingTime - safetyMargin
	if remainingTime <= safetyMargin {
		hardLimit = max(remainingTime/2, 10*time.Millisecond)
	} else if hardLimit < minSearchTime {
		hardLimit = min(remainingTime, minSearchTime)
	}
	target = min(target, hardLimit)

	return target
}
