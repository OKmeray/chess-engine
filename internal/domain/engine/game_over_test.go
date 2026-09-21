package engine

import (
	"testing"
)

func TestPosition_IsCheckmate(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"Fools mate", "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 0 3", true},
		{"Back rank mate", "4R1k1/5ppp/8/8/8/8/8/6K1 b - - 0 150", true},
		{"Stalemate is not checkmate", "k7/2Q5/8/8/8/8/8/7K b - - 0 150", false},
		{"Normal position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			if got := pos.IsCheckmate(); got != tt.want {
				t.Errorf("IsCheckmate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsStalemate(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"Stalemate", "k7/2Q5/8/8/8/8/8/7K b - - 0 150", true},
		{"Blocked Pawns Stalemate", "k7/P7/1K6/8/8/8/8/8 b - - 0 150", true},
		{"Checkmate is not stalemate", "4R1k1/5ppp/8/8/8/8/8/6K1 b - - 0 150", false},
		{"Normal position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			if got := pos.IsStalemate(); got != tt.want {
				t.Errorf("IsStalemate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsFiftyMoveRule(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"100 halfmoves white to move", "8/3k4/1b6/8/4b3/8/2K5/8 w - - 100 150", true},
		{"100 halfmoves black to move", "8/1b1k4/1b6/8/8/2K5/8/8 b - - 100 150", true},
		{"99 halfmoves", "8/1b1k4/1b6/8/8/8/2K5/8 b - - 99 149", false},
		{"Start position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			if got := pos.IsFiftyMoveRule(); got != tt.want {
				t.Errorf("IsFiftyMoveRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsInsufficientMaterial(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"K vs K", "k7/8/8/8/8/8/8/K7 w - - 0 80", true},
		{"K+N vs K", "k7/8/8/8/8/8/8/KN6 w - - 0 80", true},
		{"K+B vs K", "k7/8/8/8/8/8/8/KB6 w - - 0 80", true},
		{"K vs K+N", "kn6/8/8/8/8/8/8/K7 w - - 0 80", true},
		{"K vs K+B", "kb6/8/8/8/8/8/8/K7 w - - 0 80", true},
		{"K+B vs K+B (same color)", "k7/8/8/2b5/8/8/1B6/K7 w - - 0 80", true},
		{"K+B vs K+B (opposite color)", "k7/8/8/8/8/3b4/1B6/K7 w - - 0 80", false},
		{"K+B+B vs K+B (same color)", "k7/8/8/2b5/8/8/1B1B4/K7 w - - 0 80", true},
		{"K vs K+B+B (same color)", "k7/8/3b4/8/1b6/8/8/K7 w - - 0 80", true},
		{"K+N+N vs K", "k7/8/8/8/8/8/8/KNN5 w - - 0 80", false},
		{"K+B vs K+N", "k7/8/8/8/8/8/8/KB5n w - - 0 80", false},
		{"K+B+N vs K", "k7/8/8/8/8/8/8/KNB5 w - - 0 80", false},
		{"K+P vs K", "k7/8/8/8/8/8/P7/K7 w - - 0 80", false},
		{"K+R vs K", "k7/8/8/8/8/8/8/KR6 w - - 0 80", false},
		{"K+Q vs K", "k7/8/8/8/8/8/8/KQ6 w - - 0 80", false},
		{"Normal position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			if got := pos.IsInsufficientMaterial(); got != tt.want {
				t.Errorf("IsInsufficientMaterial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsThreeFoldRepetition(t *testing.T) {
	tests := []struct {
		name     string
		startFen string
		moves    []Move
		want     bool
	}{
		{
			name:     "No repetition",
			startFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e4"), FlagDoublePawn),
				NewMove(GetNumBySquareName("e7"), GetNumBySquareName("e5"), FlagDoublePawn),
			},
			want: false,
		},
		{
			name:     "Two-fold repetition is not enough",
			startFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e3"), FlagQuiet), // Clear start
				NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
				NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
				NewMove(GetNumBySquareName("f6"), GetNumBySquareName("g8"), FlagQuiet),
				NewMove(GetNumBySquareName("f3"), GetNumBySquareName("g1"), FlagQuiet), // Reaches position 2nd time
			},
			want: false,
		},
		{
			name:     "Three-fold repetition is detected",
			startFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e3"), FlagQuiet), // Clear start
				// 1st loop
				NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
				NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
				NewMove(GetNumBySquareName("f6"), GetNumBySquareName("g8"), FlagQuiet),
				NewMove(GetNumBySquareName("f3"), GetNumBySquareName("g1"), FlagQuiet), // Reaches position 2nd time
				// 2nd loop
				NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
				NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
				NewMove(GetNumBySquareName("f6"), GetNumBySquareName("g8"), FlagQuiet),
				NewMove(GetNumBySquareName("f3"), GetNumBySquareName("g1"), FlagQuiet), // Reaches position 3rd time
			},
			want: true,
		},
		{
			name:     "Repetition broken by pawn move",
			startFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e3"), FlagQuiet),
				// 1st loop
				NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
				NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
				NewMove(GetNumBySquareName("f6"), GetNumBySquareName("g8"), FlagQuiet),
				NewMove(GetNumBySquareName("f3"), GetNumBySquareName("g1"), FlagQuiet), // Reaches position 2nd time
				// Pawn move permanently changes state and resets HalfMoves
				NewMove(GetNumBySquareName("a7"), GetNumBySquareName("a6"), FlagQuiet),
				// Next sequence
				NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
				NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
				NewMove(GetNumBySquareName("f3"), GetNumBySquareName("g1"), FlagQuiet),
				NewMove(GetNumBySquareName("f6"), GetNumBySquareName("g8"), FlagQuiet),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.startFen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			undos := make([]UndoInfo, 0, len(tt.moves))
			for _, m := range tt.moves {
				undos = append(undos, pos.MakeMove(m))
			}

			if got := pos.IsThreeFoldRepetition(); got != tt.want {
				t.Errorf("IsThreeFoldRepetition() = %v, want %v", got, tt.want)
			}

			// Verify unmaking cleanly restores state and drops repetition count
			for i := len(tt.moves) - 1; i >= 0; i-- {
				pos.UnmakeMove(tt.moves[i], undos[i])
			}

			if got := pos.IsThreeFoldRepetition(); got != false {
				t.Errorf("IsThreeFoldRepetition() after unmake = %v, want false", got)
			}
		})
	}
}
