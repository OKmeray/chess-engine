package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

func TestMinimax_Search_Tactics(t *testing.T) {
	tests := []struct {
		name        string
		fen         string
		depth       int
		expectedSAN string
	}{
		{
			name:        "Mate in 1",
			fen:         "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 4",
			depth:       2,
			expectedSAN: "Qxf7#",
		},
		{
			name:        "Mate in 2",
			fen:         "r2q1r1k/1b2Nppp/p2p4/1p6/4P3/1P5R/1PP2PPP/3Q2K1 w - - 0 50",
			depth:       4,
			expectedSAN: "Rxh7+",
		},
		{
			name:        "Mate in 3",
			fen:         "1k1r4/pp1b1R2/3q2pp/4p3/2B5/4Q3/PPP2B2/2K5 b - - 0 50",
			depth:       6,
			expectedSAN: "Qd1+",
		},
		{
			name:        "Knight fork in 2 moves",
			fen:         "6k1/5r1p/p2N4/nppP2q1/2P5/1P2N3/PQ5P/7K w - - 0 30",
			depth:       5,
			expectedSAN: "Qh8+",
		},
		{
			name:        "Pawn promotion",
			fen:         "6k1/P4pp1/7p/3q4/7n/8/R4PPP/7K w - - 0 50",
			depth:       3,
			expectedSAN: "a8=Q",
		},
		{
			name:        "Prefer faster mate",
			fen:         "4r1k1/pp3p1p/2p2Qp1/6N1/8/8/PPP2PPP/4R1K1 w - - 0 50",
			depth:       3,
			expectedSAN: "Rxe8#",
		},
		{
			name:        "Force stalemate",
			fen:         "8/8/8/8/8/3K1ppp/5pbq/5krn w - - 0 150",
			depth:       3,
			expectedSAN: "Kd2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := engine.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			searcher := NewMinimax(tt.depth)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			bestMove, _ := searcher.Search(ctx, pos)

			expectedMove, err := engine.ParseSAN(pos, tt.expectedSAN)
			if err != nil {
				t.Fatalf("Failed to parse expected SAN %s: %v", tt.expectedSAN, err)
			}

			if bestMove != expectedMove {
				t.Errorf("Search() move = %v, want %v", bestMove, expectedMove)
			}
		})
	}
}

func TestMinimax_Search_GameOver(t *testing.T) {
	tests := []struct {
		name          string
		fen           string
		expectedScore float32
	}{
		{
			name:          "Stalemate",
			fen:           "k7/P7/1K6/8/8/8/8/8 b - - 0 1",
			expectedScore: 0.0,
		},
		{
			name:          "Checkmate",
			fen:           "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3",
			expectedScore: -10002.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := engine.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			searcher := NewMinimax(2)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			bestMove, score := searcher.Search(ctx, pos)

			if bestMove != 0 {
				t.Errorf("Search() move = %v, want %v", bestMove, 0)
			}
			if score != tt.expectedScore {
				t.Errorf("Search() score = %f, want %f", score, tt.expectedScore)
			}
		})
	}
}

type mockContext struct {
	context.Context
	checks      int
	cancelAfter int
}

func (m *mockContext) Err() error {
	m.checks++
	if m.checks > m.cancelAfter {
		return context.Canceled
	}
	return nil
}

func TestMinimax_Search_Timeout(t *testing.T) {
	pos, err := engine.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	searcher := NewMinimax(20)

	// we use a deterministic mock context that cancels itself
	// after Err() is called 5 times to imitate timeout cancellation
	ctx := &mockContext{Context: context.Background(), cancelAfter: 5}

	bestMove, _ := searcher.Search(ctx, pos)

	if bestMove == 0 {
		t.Errorf("Search() move = 0, want valid move from previous depth")
	}
}

func TestMinimax_Search_ImmediateCancel(t *testing.T) {
	pos, err := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	searcher := NewMinimax(5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	bestMove, _ := searcher.Search(ctx, pos)
	if bestMove != 0 {
		t.Errorf("Search() move = %v, want 0", bestMove)
	}
}
