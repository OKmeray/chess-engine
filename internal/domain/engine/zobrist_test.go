package engine

import (
	"testing"
)

func TestZobristIncrementalUpdate(t *testing.T) {
	tests := []struct {
		name string
		fen  string
	}{
		{
			name: "Starting Position",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "Kiwipete (Complex Middlegame)",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		},
		{
			name: "En Passant Capture Possible",
			fen:  "8/8/8/8/5pP1/8/8/4K2k b - g3 0 1",
		},
		{
			name: "Castling Destroyed",
			fen:  "r3k2r/p6p/8/8/8/8/P6P/R3K2R w KQkq - 0 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			initialHash := pos.Hash
			computedInitial := pos.ComputeHash()
			if initialHash != computedInitial {
				t.Fatalf("ComputeHash() = %x, want %x", computedInitial, initialHash)
			}

			var movesBuf [256]Move
			moves := pos.GenerateMoves(movesBuf[:0])

			for _, move := range moves {
				undo := pos.MakeMove(move)

				incrementalHash := pos.Hash
				computedHash := pos.ComputeHash()

				if incrementalHash != computedHash {
					t.Errorf("MakeMove(%v) incremental hash = %x, want %x", move, incrementalHash, computedHash)
				}

				pos.UnmakeMove(move, undo)

				restoredHash := pos.Hash
				if restoredHash != initialHash {
					t.Errorf("UnmakeMove(%v) restored hash = %x, want %x", move, restoredHash, initialHash)
				}
			}
		})
	}
}
