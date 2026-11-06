package usecase

import (
	"context"
	"fmt"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// EngineType represents the algorithm used for move search calculations.
type EngineType string

const (
	// EngineTypeMCTS specifies the Monte Carlo Tree Search algorithm.
	EngineTypeMCTS EngineType = "mcts"
	// EngineTypeMinimax specifies the Alpha-Beta Minimax search algorithm.
	EngineTypeMinimax EngineType = "minimax"
)

// Searcher defines how a chess position is searched for the best move.
type Searcher interface {
	Search(ctx context.Context, pos *engine.Position) (engine.Move, error)
}

var (
	_ Searcher = (*Minimax)(nil)
	_ Searcher = (*MCTSSearcher)(nil)
)

// EvaluatorProvider retrieves loaded neural network evaluators by model ID.
type EvaluatorProvider interface {
	Get(ctx context.Context, modelID int) (Evaluator, error)
}

// GameCreator instantiates game sessions configured with the requested search algorithm.
type GameCreator struct {
	evaluators   EvaluatorProvider
	maxNodes     int
	numWorkers   int
	minimaxDepth int
}

// NewGameCreator initializes a GameCreator with engine dependencies and limits.
func NewGameCreator(evaluators EvaluatorProvider, maxNodes, numWorkers, minimaxDepth int) *GameCreator {
	return &GameCreator{
		evaluators:   evaluators,
		maxNodes:     maxNodes,
		numWorkers:   numWorkers,
		minimaxDepth: minimaxDepth,
	}
}

// CreateGame creates a MoveFinder game session for the given engineType and botID.
func (c *GameCreator) CreateGame(ctx context.Context, engineType EngineType, botID int) (MoveFinder, error) {
	switch engineType {
	case EngineTypeMCTS:
		if c == nil || c.evaluators == nil {
			return nil, ErrNilEvaluatorProvider
		}
		eval, err := c.evaluators.Get(ctx, botID)
		if err != nil {
			return nil, err
		}
		return NewGame(NewMCTSSearcher(eval, c.maxNodes, c.numWorkers))
	case EngineTypeMinimax:
		return NewGame(NewMinimax(c.minimaxDepth))
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedEngineType, engineType)
	}
}
