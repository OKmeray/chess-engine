package engine

import "testing"

func TestNewPosition(t *testing.T) {
	pos := NewPosition()
	if pos.SideToMove != White {
		t.Errorf("NewPosition() SideToMove = %v, want %v", pos.SideToMove, White)
	}
	if pos.EnPassantSquare != NoSquare {
		t.Errorf("NewPosition() EnPassantSquare = %v, want %v", pos.EnPassantSquare, NoSquare)
	}
	if pos.HalfMoves != 0 {
		t.Errorf("NewPosition() HalfMoves = %v, want %v", pos.HalfMoves, 0)
	}
	if pos.CurrentTurn != 1 {
		t.Errorf("NewPosition() CurrentTurn = %v, want %v", pos.CurrentTurn, 1)
	}

	expectedCastlingRights := CastlingRights(^WhiteShort | ^WhiteLong | ^BlackShort | ^BlackLong)
	if pos.CastlingRights != expectedCastlingRights {
		t.Errorf("NewPosition() CastlingRights = %v, want %v", pos.CastlingRights, expectedCastlingRights)
	}

	// Check that bitboards are empty
	for i, bb := range pos.Pieces {
		if bb != 0 {
			t.Errorf("NewPosition() Pieces[%d] = %v, want %v", i, bb, 0)
		}
	}
	for i, bb := range pos.ByColor {
		if bb != 0 {
			t.Errorf("NewPosition() ByColor[%d] = %v, want %v", i, bb, 0)
		}
	}
}

func TestPosition_MakeUnmakeMove(t *testing.T) {
	tests := []struct {
		name        string
		startFen    string
		move        Move
		expectedFen string
	}{
		// Pawn Pushes
		{
			name:        "White pawn single push",
			startFen:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e3"), FlagQuiet),
			expectedFen: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:        "Black pawn single push",
			startFen:    "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			move:        NewMove(GetNumBySquareName("e7"), GetNumBySquareName("e6"), FlagQuiet),
			expectedFen: "rnbqkbnr/pppp1ppp/4p3/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
		},
		{
			name:        "White double pawn push",
			startFen:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e4"), FlagDoublePawn),
			expectedFen: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name:        "Black double pawn push",
			startFen:    "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			move:        NewMove(GetNumBySquareName("e7"), GetNumBySquareName("e5"), FlagDoublePawn),
			expectedFen: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},

		// Pawn Captures
		{
			name:        "White pawn captures black pawn",
			startFen:    "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			move:        NewMove(GetNumBySquareName("e4"), GetNumBySquareName("d5"), FlagCapture),
			expectedFen: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			name:        "Black pawn captures white pawn",
			startFen:    "rnbqkbnr/ppp1pppp/8/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			move:        NewMove(GetNumBySquareName("d5"), GetNumBySquareName("e4"), FlagCapture),
			expectedFen: "rnbqkbnr/ppp1pppp/8/8/4p3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3",
		},
		{
			name:        "White pawn captures en passant",
			startFen:    "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
			move:        NewMove(GetNumBySquareName("e5"), GetNumBySquareName("f6"), FlagEnPassant),
			expectedFen: "rnbqkbnr/ppp1p1pp/5P2/3p4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		{
			name:        "Black pawn captures en passant",
			startFen:    "rnbqkbnr/ppp1pppp/8/4P3/2Pp4/8/PP1P1PPP/RNBQKBNR b KQkq c3 0 3",
			move:        NewMove(GetNumBySquareName("d4"), GetNumBySquareName("c3"), FlagEnPassant),
			expectedFen: "rnbqkbnr/ppp1pppp/8/4P3/8/2p5/PP1P1PPP/RNBQKBNR w KQkq - 0 4",
		},

		// Pawn Promotions
		{
			name:        "White pawn promotion to queen",
			startFen:    "4k3/1P6/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("b8"), FlagPromoQueen),
			expectedFen: "1Q2k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "White pawn promotion to rook",
			startFen:    "4k3/1P6/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("b8"), FlagPromoRook),
			expectedFen: "1R2k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "White pawn promotion to bishop",
			startFen:    "4k3/1P6/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("b8"), FlagPromoBishop),
			expectedFen: "1B2k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "White pawn promotion to knight",
			startFen:    "4k3/1P6/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("b8"), FlagPromoKnight),
			expectedFen: "1N2k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "Black pawn promotion to queen",
			startFen:    "4k3/8/8/8/8/8/1p6/4K3 b - - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("b1"), FlagPromoQueen),
			expectedFen: "4k3/8/8/8/8/8/8/1q2K3 w - - 0 2",
		},
		{
			name:        "Black pawn promotion to rook",
			startFen:    "4k3/8/8/8/8/8/1p6/4K3 b - - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("b1"), FlagPromoRook),
			expectedFen: "4k3/8/8/8/8/8/8/1r2K3 w - - 0 2",
		},
		{
			name:        "Black pawn promotion to bishop",
			startFen:    "4k3/8/8/8/8/8/1p6/4K3 b - - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("b1"), FlagPromoBishop),
			expectedFen: "4k3/8/8/8/8/8/8/1b2K3 w - - 0 2",
		},
		{
			name:        "Black pawn promotion to knight",
			startFen:    "4k3/8/8/8/8/8/1p6/4K3 b - - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("b1"), FlagPromoKnight),
			expectedFen: "4k3/8/8/8/8/8/8/1n2K3 w - - 0 2",
		},
		{
			name:        "White pawn captures enemy rook and promotes to queen",
			startFen:    "2r1k3/1P6/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("c8"), FlagPromoQueenCapture),
			expectedFen: "2Q1k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "Black pawn captures enemy knight and promotes to bishop",
			startFen:    "4k3/8/8/8/8/8/1p6/2N1K3 b - - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("c1"), FlagPromoBishopCapture),
			expectedFen: "4k3/8/8/8/8/8/8/2b1K3 w - - 0 2",
		},

		// Castling
		{
			name:        "White short castling",
			startFen:    "r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			move:        NewMove(GetNumBySquareName("e1"), GetNumBySquareName("g1"), FlagKingCastle),
			expectedFen: "r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
		},
		{
			name:        "White long castling",
			startFen:    "r3kbnr/pppqpppp/2n5/3p1b2/3P1B2/2N5/PPPQPPPP/R3KBNR w KQkq - 4 5",
			move:        NewMove(GetNumBySquareName("e1"), GetNumBySquareName("c1"), FlagQueenCastle),
			expectedFen: "r3kbnr/pppqpppp/2n5/3p1b2/3P1B2/2N5/PPPQPPPP/2KR1BNR b kq - 5 5",
		},
		{
			name:        "Black short castling",
			startFen:    "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQ1RK1 b kq - 0 5",
			move:        NewMove(GetNumBySquareName("e8"), GetNumBySquareName("g8"), FlagKingCastle),
			expectedFen: "r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQ1RK1 w - - 1 6",
		},
		{
			name:        "Black long castling",
			startFen:    "r3kbnr/pppqpppp/2n5/3p1b2/3P1B2/2N5/PPPQPPPP/2KR1BNR b kq - 5 5",
			move:        NewMove(GetNumBySquareName("e8"), GetNumBySquareName("c8"), FlagQueenCastle),
			expectedFen: "2kr1bnr/pppqpppp/2n5/3p1b2/3P1B2/2N5/PPPQPPPP/2KR1BNR w - - 6 6",
		},

		// Castling Rights Revocation
		{
			name:        "White king move revokes all castling rights",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("e1"), GetNumBySquareName("e2"), FlagQuiet),
			expectedFen: "r3k2r/8/8/8/8/8/4K3/R6R b kq - 1 1",
		},
		{
			name:        "White kingside rook move revokes kingside castling",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("h1"), GetNumBySquareName("h2"), FlagQuiet),
			expectedFen: "r3k2r/8/8/8/8/8/7R/R3K3 b Qkq - 1 1",
		},
		{
			name:        "White queenside rook move revokes queenside castling",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("a1"), GetNumBySquareName("a2"), FlagQuiet),
			expectedFen: "r3k2r/8/8/8/8/8/R7/4K2R b Kkq - 1 1",
		},
		{
			name:        "Black kingside rook move revokes kingside castling",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("h8"), GetNumBySquareName("h7"), FlagQuiet),
			expectedFen: "r3k3/7r/8/8/8/8/8/R3K2R w KQq - 1 2",
		},
		{
			name:        "Black queenside rook move revokes queenside castling",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("a8"), GetNumBySquareName("a7"), FlagQuiet),
			expectedFen: "4k2r/r7/8/8/8/8/8/R3K2R w KQk - 1 2",
		},
		{
			name:        "Black rook capture revokes castling rights",
			startFen:    "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("h8"), GetNumBySquareName("h1"), FlagCapture),
			expectedFen: "r3k3/8/8/8/8/8/8/R3K2r w Qq - 0 2",
		},
		{
			name:        "White captures black queenside rook revokes black long castling",
			startFen:    "r3k2r/1B6/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("b7"), GetNumBySquareName("a8"), FlagCapture),
			expectedFen: "B3k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1",
		},
		{
			name:        "White captures black kingside rook revokes black short castling",
			startFen:    "r3k2r/6B1/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("g7"), GetNumBySquareName("h8"), FlagCapture),
			expectedFen: "r3k2B/8/8/8/8/8/8/R3K2R b KQq - 0 1",
		},
		{
			name:        "Black captures white kingside rook revokes white short castling",
			startFen:    "r3k2r/8/8/8/8/8/6b1/R3K2R b KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("g2"), GetNumBySquareName("h1"), FlagCapture),
			expectedFen: "r3k2r/8/8/8/8/8/8/R3K2b w Qkq - 0 2",
		},
		{
			name:        "Black captures white queenside rook revokes white long castling",
			startFen:    "r3k2r/8/8/8/8/8/1b6/R3K2R b KQkq - 0 1",
			move:        NewMove(GetNumBySquareName("b2"), GetNumBySquareName("a1"), FlagCapture),
			expectedFen: "r3k2r/8/8/8/8/8/8/b3K2R w Kkq - 0 2",
		},

		// Special Cases
		{
			name:        "Knight moves to deliver double check (discovered attack)",
			startFen:    "k3N2R/8/8/8/8/8/8/4K3 w - - 0 1",
			move:        NewMove(GetNumBySquareName("e8"), GetNumBySquareName("c7"), FlagQuiet),
			expectedFen: "k6R/2N5/8/8/8/8/8/4K3 b - - 1 1",
		},
		{
			name:        "King evades check by moving to empty square",
			startFen:    "k7/8/8/8/8/8/8/R3K3 b Q - 0 1",
			move:        NewMove(GetNumBySquareName("a8"), GetNumBySquareName("b7"), FlagQuiet),
			expectedFen: "8/1k6/8/8/8/8/8/R3K3 w Q - 1 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.startFen)
			if err != nil {
				t.Fatalf("failed to parse start FEN: %v", err)
			}

			undo := pos.MakeMove(tt.move)

			gotFen := pos.FEN()
			if gotFen != tt.expectedFen {
				t.Errorf("MakeMove() FEN = %q, want %q", gotFen, tt.expectedFen)
			}

			pos.UnmakeMove(tt.move, undo)

			gotStartFen := pos.FEN()
			if gotStartFen != tt.startFen {
				t.Errorf("UnmakeMove() FEN = %q, want %q", gotStartFen, tt.startFen)
			}
		})
	}
}

func TestPosition_IsSquareAttacked(t *testing.T) {
	tests := []struct {
		name          string
		fen           string
		square        string
		attackerColor PieceColor
		want          bool
	}{
		{
			name:          "White pawn e4 attacks d5",
			fen:           "8/8/8/8/4P3/8/8/K6k w - - 0 1",
			square:        "d5",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White pawn e4 attacks f5",
			fen:           "8/8/8/8/4P3/8/8/K6k w - - 0 1",
			square:        "f5",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White pawn e4 does not attack e5",
			fen:           "8/8/8/8/4P3/8/8/K6k w - - 0 1",
			square:        "e5",
			attackerColor: White,
			want:          false,
		},
		{
			name:          "Black pawn e5 attacks f4",
			fen:           "8/8/8/4p3/8/8/8/K6k b - - 0 1",
			square:        "f4",
			attackerColor: Black,
			want:          true,
		},
		{
			name:          "Black pawn e5 attacks d4",
			fen:           "8/8/8/4p3/8/8/8/K6k b - - 0 1",
			square:        "d4",
			attackerColor: Black,
			want:          true,
		},
		{
			name:          "Black pawn e5 does not attacks e4",
			fen:           "8/8/8/4p3/8/8/8/K6k b - - 0 1",
			square:        "e4",
			attackerColor: Black,
			want:          false,
		},
		{
			name:          "White knight f3 attacks d4",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "d4",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White bishop c4 attacks f7",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "f7",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White rook e1 attacks e7",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "e7",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White queen d4 attacks g7",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "g7",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White king g1 attacks f1",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "f1",
			attackerColor: White,
			want:          true,
		},
		{
			name:          "White rook c1 attack on c6 blocked by c4",
			fen:           "r4rk1/1pq1bppp/2n5/2p5/2BQ4/5N1P/PP3PP1/2R1R1K1 w - - 0 1",
			square:        "c6",
			attackerColor: White,
			want:          false,
		},
		{
			name:          "Pinned black knight c6 attacks e5",
			fen:           "4k3/8/2n5/1B6/8/8/8/4K3 b - - 0 1",
			square:        "e5",
			attackerColor: Black,
			want:          true,
		},
		{
			name:          "Black pawn d5 attacks e4",
			fen:           "4k3/8/8/3pP3/8/8/8/4K3 b - - 0 1",
			square:        "e4",
			attackerColor: Black,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			sq := GetNumBySquareName(tt.square)
			got := pos.IsSquareAttacked(sq, tt.attackerColor)
			if got != tt.want {
				t.Errorf("IsSquareAttacked(%v, %v) = %v, want %v", tt.square, tt.attackerColor, got, tt.want)
			}
		})
	}
}

func TestPosition_GetKingSq(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		color PieceColor
		want  string
	}{
		{
			name:  "White king on e1",
			fen:   "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
			color: White,
			want:  "e1",
		},
		{
			name:  "Black king on e8",
			fen:   "4k3/8/8/8/8/8/8/4K3 b - - 0 1",
			color: Black,
			want:  "e8",
		},
		{
			name:  "White king on h1",
			fen:   "k7/8/8/8/8/8/8/7K w - - 0 1",
			color: White,
			want:  "h1",
		},
		{
			name:  "Black king on a8",
			fen:   "k7/8/8/8/8/8/8/7K w - - 0 1",
			color: Black,
			want:  "a8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			wantSq := GetNumBySquareName(tt.want)
			got := pos.GetKingSq(tt.color)
			if got != wantSq {
				t.Errorf("GetKingSq(%v) = %v (%v), want %v (%v)", tt.color, GetSquareNameByNum(got), got, tt.want, wantSq)
			}
		})
	}
}

func TestPosition_HalfMoves(t *testing.T) {
	startFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(startFen)
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	tests := []struct {
		name          string
		move          Move
		wantHalfMoves int
	}{
		{
			name:          "White Knight makes quite move",
			move:          NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet),
			wantHalfMoves: 1,
		},
		{
			name:          "Black Knight makes quite move",
			move:          NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet),
			wantHalfMoves: 2,
		},
		{
			name:          "White Knight makes quite move",
			move:          NewMove(GetNumBySquareName("b1"), GetNumBySquareName("c3"), FlagQuiet),
			wantHalfMoves: 3,
		},
		{
			name:          "Black Knight makes quite move",
			move:          NewMove(GetNumBySquareName("b8"), GetNumBySquareName("c6"), FlagQuiet),
			wantHalfMoves: 4,
		},
		{
			name:          "White Double Pawn Push",
			move:          NewMove(GetNumBySquareName("e2"), GetNumBySquareName("e4"), FlagDoublePawn),
			wantHalfMoves: 0,
		},
		{
			name:          "Black Double Pawn Push",
			move:          NewMove(GetNumBySquareName("e7"), GetNumBySquareName("e5"), FlagDoublePawn),
			wantHalfMoves: 0,
		},
		{
			name:          "White Bishop makes quite move",
			move:          NewMove(GetNumBySquareName("f1"), GetNumBySquareName("c4"), FlagQuiet),
			wantHalfMoves: 1,
		},
		{
			name:          "Black Bishop makes quite move",
			move:          NewMove(GetNumBySquareName("f8"), GetNumBySquareName("c5"), FlagQuiet),
			wantHalfMoves: 2,
		},
		{
			name:          "Castling (White King side)",
			move:          NewMove(GetNumBySquareName("e1"), GetNumBySquareName("g1"), FlagKingCastle),
			wantHalfMoves: 3,
		},
		{
			name:          "Black Knight makes quite move",
			move:          NewMove(GetNumBySquareName("c6"), GetNumBySquareName("d4"), FlagQuiet),
			wantHalfMoves: 4,
		},
		{
			name:          "White Knight captures Black Knight",
			move:          NewMove(GetNumBySquareName("f3"), GetNumBySquareName("d4"), FlagCapture),
			wantHalfMoves: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos.MakeMove(tt.move)
			if pos.HalfMoves != tt.wantHalfMoves {
				t.Errorf("MakeMove() HalfMoves = %d, want %d", pos.HalfMoves, tt.wantHalfMoves)
			}
		})
	}
}

func TestPosition_RepetitionHistory(t *testing.T) {
	startFen := "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"
	pos, err := ParseFEN(startFen)
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	// moves for the test
	// 5 quiet moves, then pawn move which resets the history, then 2 more
	// quiet moves, then we undo 3 moves and test if the repetition history
	// contains correct hashes for the first 5 moves
	m1 := NewMove(GetNumBySquareName("g1"), GetNumBySquareName("f3"), FlagQuiet)
	m2 := NewMove(GetNumBySquareName("b8"), GetNumBySquareName("c6"), FlagQuiet)
	m3 := NewMove(GetNumBySquareName("f1"), GetNumBySquareName("c4"), FlagQuiet)
	m4 := NewMove(GetNumBySquareName("f8"), GetNumBySquareName("c5"), FlagQuiet)
	m5 := NewMove(GetNumBySquareName("b1"), GetNumBySquareName("c3"), FlagQuiet)
	m6 := NewMove(GetNumBySquareName("d7"), GetNumBySquareName("d6"), FlagQuiet)
	m7 := NewMove(GetNumBySquareName("e1"), GetNumBySquareName("g1"), FlagKingCastle)
	m8 := NewMove(GetNumBySquareName("g8"), GetNumBySquareName("f6"), FlagQuiet)

	moves := []Move{m1, m2, m3, m4, m5, m6, m7, m8}
	undos := make([]UndoInfo, 8)
	expectedHashes := make([]uint64, 8)

	for i, m := range moves {
		expectedHashes[i] = pos.Hash
		undos[i] = pos.MakeMove(m)
	}

	// Now unmake 3 moves (m8, m7, m6) to get back to the state after move 5
	for i := 7; i >= 5; i-- {
		pos.UnmakeMove(moves[i], undos[i])
	}

	// Verify the RepetitionHistory is restored correctly for the first 5 moves
	for i := 0; i < 5; i++ {
		if pos.RepetitionHistory[i] != expectedHashes[i] {
			t.Errorf("after unmaking back to move 5, RepetitionHistory[%d] = %d, want %d", i, pos.RepetitionHistory[i], expectedHashes[i])
		}
	}

	// Verify unmaking the rest works
	for i := 4; i >= 0; i-- {
		pos.UnmakeMove(moves[i], undos[i])
	}

	if pos.FEN() != startFen {
		t.Errorf("FEN after fully unmaking = %q, want %q", pos.FEN(), startFen)
	}
}
