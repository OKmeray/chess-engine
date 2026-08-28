package usecase

import (
	"context"
	"math"
	"sort"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// Minimax represents the Minimax search algorithm with Alpha-Beta pruning.
type Minimax struct {
	maxDepth      int
	nodesSearched int
}

// NewMinimax creates a new Minimax searcher.
func NewMinimax(maxDepth int) *Minimax {
	return &Minimax{
		maxDepth: maxDepth,
	}
}

// Search finds the best move in the given position using iterative deepening.
func (m *Minimax) Search(ctx context.Context, pos *engine.Position) (engine.Move, float32) {
	var bestMove engine.Move
	var bestScore float32
	m.nodesSearched = 0

	// Iterative Deepening
	for depth := 1; depth <= m.maxDepth; depth++ {
		move, score, aborted := m.searchDepth(ctx, pos, depth)

		if aborted {
			// If we aborted on depth 1, we just have to return whatever we found.
			// Otherwise, discard the incomplete depth and return the previous depth's result.
			if depth == 1 {
				bestMove = move
				bestScore = score
			}
			break
		}

		bestMove = move
		bestScore = score

		// If we found a forced checkmate, end the search
		if bestScore >= 5000.0 {
			break
		}
	}

	return bestMove, bestScore
}

// searchDepth finds the best move by initiating an alpha-beta search to the given depth.
func (m *Minimax) searchDepth(ctx context.Context, pos *engine.Position, depth int) (engine.Move, float32, bool) {
	var bestMove engine.Move
	bestScore := float32(math.Inf(-1))

	movesBuf := make([]engine.Move, 0, 256)
	moves := pos.GenerateMoves(movesBuf)

	if len(moves) == 0 {
		if pos.IsSquareAttacked(pos.GetKingSq(pos.SideToMove), 6-pos.SideToMove) {
			return 0, -10000.0 - float32(depth), false
		}
		return 0, 0, false
	}

	sortMoves(pos, moves)

	alpha := float32(math.Inf(-1))
	beta := float32(math.Inf(1))

	for _, move := range moves {
		undo := pos.MakeMove(move)
		score, aborted := m.alphaBeta(ctx, pos, depth-1, -beta, -alpha)
		pos.UnmakeMove(move, undo)
		score = -score

		if aborted {
			return bestMove, bestScore, true
		}

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
		}
	}

	return bestMove, bestScore, false
}

// alphaBeta is the recursive negamax function with alpha-beta pruning.
func (m *Minimax) alphaBeta(ctx context.Context, pos *engine.Position, depth int, alpha, beta float32) (float32, bool) {
	m.nodesSearched++

	// Check for timeout
	if (m.nodesSearched&2047 == 0 || m.nodesSearched == 1) && ctx.Err() != nil {
		return 0, true
	}

	if depth == 0 {
		score := evaluate(pos)
		if pos.SideToMove == engine.Black {
			score = -score
		}
		return score, false
	}

	var movesBuf [256]engine.Move
	moves := pos.GenerateMoves(movesBuf[:0])

	if len(moves) == 0 {
		if pos.IsSquareAttacked(pos.GetKingSq(pos.SideToMove), 6-pos.SideToMove) {
			return -10000.0 - float32(depth), false // Checkmate
		}
		return 0.0, false // Stalemate
	}

	sortMoves(pos, moves)

	bestScore := float32(math.Inf(-1))
	for _, move := range moves {
		undo := pos.MakeMove(move)
		score, aborted := m.alphaBeta(ctx, pos, depth-1, -beta, -alpha)
		pos.UnmakeMove(move, undo)
		score = -score

		if aborted {
			return 0, true
		}

		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	return bestScore, false
}

var piecePrices = [6]int{
	int(engine.PawnPrice),
	int(engine.KnightPrice),
	int(engine.BishopPrice),
	int(engine.RookPrice),
	int(engine.QueenPrice),
	int(engine.KingPrice),
}

// sortMoves sorts moves in place using MVV-LVA and Check detection.
func sortMoves(pos *engine.Position, moves []engine.Move) {
	type moveScore struct {
		move  engine.Move
		score int
	}
	moveScores := make([]moveScore, len(moves))

	enemyKingColor := 6 - pos.SideToMove
	enemyKingSq := pos.GetKingSq(enemyKingColor)

	for i, move := range moves {
		score := 0
		flags := move.Flags()

		if flags&engine.FlagCapture != 0 || flags == engine.FlagEnPassant {
			victimType := engine.Pawn
			if flags != engine.FlagEnPassant {
				victimType, _ = pos.GetPieceAndColorBySquare(move.To())
			}
			attackerType, _ := pos.GetPieceAndColorBySquare(move.From())

			if victimType != engine.None && attackerType != engine.None {
				score += 20000 + 10*piecePrices[victimType] - piecePrices[attackerType]
			}
		}

		if flags >= engine.FlagPromoKnight && flags <= engine.FlagPromoQueenCapture {
			score += 30000
		}

		undo := pos.MakeMove(move)
		if pos.IsSquareAttacked(enemyKingSq, 6-pos.SideToMove) {
			score += 10000
		}
		pos.UnmakeMove(move, undo)

		moveScores[i] = moveScore{move: move, score: score}
	}

	sort.SliceStable(moveScores, func(i, j int) bool {
		return moveScores[i].score > moveScores[j].score
	})

	for i, ms := range moveScores {
		moves[i] = ms.move
	}
}

// evaluate is a basic material-based heuristic evaluation.
func evaluate(pos *engine.Position) float32 {
	var score float32
	for i := 0; i < 64; i++ {
		pt, color := pos.GetPieceAndColorBySquare(engine.Square(i))
		if pt != engine.None {
			val := float32(piecePrices[pt])
			if color == engine.White {
				score += val
			} else {
				score -= val
			}
		}
	}
	return score
}
