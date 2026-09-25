package usecase

import (
	"context"
	"math"
	"sort"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// TTFlag represents the type of score bound stored in the Transposition Table.
type TTFlag uint8

const (
	// TTExact indicates the stored score is an exact minimax value.
	TTExact TTFlag = iota
	// TTAlpha indicates the stored score is an upper bound.
	TTAlpha
	// TTBeta indicates the stored score is a lower bound.
	TTBeta
)

// TTEntry represents a cached search result in the Transposition Table.
type TTEntry struct {
	Hash  uint64
	Move  engine.Move
	Score float32
	Depth int8
	Flag  TTFlag
}

// Minimax represents the Minimax search algorithm with Alpha-Beta pruning.
type Minimax struct {
	maxDepth      int
	nodesSearched int
	tt            []TTEntry
}

// NewMinimax creates a new Minimax searcher.
func NewMinimax(maxDepth int) *Minimax {
	return &Minimax{
		maxDepth: maxDepth,
		tt:       make([]TTEntry, 8*1024*1024), // ~8*24MB Transposition Table (8M entries)
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

	if pos.IsFiftyMoveRule() || pos.IsInsufficientMaterial() || pos.IsThreeFoldRepetition() {
		return 0, 0, false // Draw
	}

	movesBuf := make([]engine.Move, 0, 256)
	moves := pos.GenerateMoves(movesBuf)

	if len(moves) == 0 {
		if pos.IsSquareAttacked(pos.GetKingSq(pos.SideToMove), 6-pos.SideToMove) {
			return 0, -10000.0 - float32(depth), false
		}
		return 0, 0, false
	}

	// Check TT for root move ordering
	var ttMove engine.Move
	ttIndex := pos.Hash & uint64(len(m.tt)-1)
	if m.tt[ttIndex].Hash == pos.Hash {
		ttMove = m.tt[ttIndex].Move
	}

	sortMoves(pos, moves, ttMove)

	alpha := float32(math.Inf(-1))
	beta := float32(math.Inf(1))
	originalAlpha := alpha

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

	flag := TTAlpha
	if bestScore > originalAlpha {
		flag = TTExact
	}
	m.tt[ttIndex] = TTEntry{
		Hash:  pos.Hash,
		Move:  bestMove,
		Score: bestScore,
		Depth: int8(depth),
		Flag:  flag,
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

	if pos.IsFiftyMoveRule() || pos.IsInsufficientMaterial() || pos.IsThreeFoldRepetition() {
		return 0.0, false // Draw
	}

	// Check TT
	originalAlpha := alpha
	ttIndex := pos.Hash & uint64(len(m.tt)-1)
	entry := m.tt[ttIndex]
	var ttMove engine.Move

	if entry.Hash == pos.Hash {
		ttMove = entry.Move
		if int(entry.Depth) >= depth {
			if entry.Flag == TTExact {
				return entry.Score, false
			}
			if entry.Flag == TTAlpha && entry.Score <= alpha {
				return alpha, false
			}
			if entry.Flag == TTBeta && entry.Score >= beta {
				return beta, false
			}
		}
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

	sortMoves(pos, moves, ttMove)

	bestScore := float32(math.Inf(-1))
	var bestMove engine.Move

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
			bestMove = move
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	flag := TTAlpha
	if bestScore >= beta {
		flag = TTBeta
	} else if bestScore > originalAlpha {
		flag = TTExact
	}

	m.tt[ttIndex] = TTEntry{
		Hash:  pos.Hash,
		Move:  bestMove,
		Score: bestScore,
		Depth: int8(depth),
		Flag:  flag,
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

// sortMoves sorts moves in place using MVV-LVA, Check detection, and TT Move.
func sortMoves(pos *engine.Position, moves []engine.Move, ttMove engine.Move) {
	type moveScore struct {
		move  engine.Move
		score int
	}
	moveScores := make([]moveScore, len(moves))

	enemyKingColor := 6 - pos.SideToMove
	enemyKingSq := pos.GetKingSq(enemyKingColor)

	for i, move := range moves {
		score := 0

		if move == ttMove {
			score += 1000000 // TT Move gets absolute highest priority
		}

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
