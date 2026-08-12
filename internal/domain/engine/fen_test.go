package engine

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseFEN_Errors(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		wantErr error
	}{
		{
			name:    "Valid starting FEN",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			wantErr: nil,
		},
		{
			name:    "Empty FEN",
			fen:     "",
			wantErr: ErrEmptyFEN,
		},
		{
			name:    "Missing fields",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0", // missing fullmove clock
			wantErr: ErrInvalidFENFormat,
		},
		{
			name:    "Invalid piece placement length",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR/8 w KQkq - 0 1",
			wantErr: ErrInvalidFENLength,
		},
		{
			name:    "Invalid piece character",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNZ w KQkq - 0 1",
			wantErr: ErrInvalidFENPiece,
		},
		{
			name:    "Invalid side to move",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1",
			wantErr: ErrInvalidFENSide,
		},
		{
			name:    "Invalid castling right",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkqX - 0 1",
			wantErr: ErrInvalidFENCastling,
		},
		{
			name:    "Invalid en passant square",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq x9 0 1",
			wantErr: ErrInvalidFENEnPassant,
		},
		{
			name:    "Invalid halfmove clock",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - x 1",
			wantErr: strconv.ErrSyntax,
		},
		{
			name:    "Invalid fullmove number",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 x",
			wantErr: strconv.ErrSyntax,
		},
		{
			name:    "Missing king",
			fen:     "rnbq1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1BNR w KQkq - 0 1", // no kings on the board
			wantErr: ErrMissingKing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFEN(tt.fen)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseFEN(%q) error = %v, wantErr %v", tt.fen, err, tt.wantErr)
			}
		})
	}
}

func TestParseFEN_State(t *testing.T) {
	// Helper to create exact internal bitboards using readable square names
	makeBitboard := func(squares ...string) Bitboard {
		var b Bitboard
		for _, sq := range squares {
			b |= Bitboard(1) << GetNumBySquareName(sq)
		}
		return b
	}

	// Helper to verify all pieces of a color in certain position
	verifyPieces := func(t *testing.T, pos *Position, color PieceColor, pawns, knights, bishops, rooks, queens, kings Bitboard) {
		t.Helper()
		if got := pos.Pieces[int(Pawn)+int(color)]; got != pawns {
			t.Errorf("Color %v Pawns = 0x%X, want 0x%X", color, got, pawns)
		}
		if got := pos.Pieces[int(Knight)+int(color)]; got != knights {
			t.Errorf("Color %v Knights = 0x%X, want 0x%X", color, got, knights)
		}
		if got := pos.Pieces[int(Bishop)+int(color)]; got != bishops {
			t.Errorf("Color %v Bishops = 0x%X, want 0x%X", color, got, bishops)
		}
		if got := pos.Pieces[int(Rook)+int(color)]; got != rooks {
			t.Errorf("Color %v Rooks = 0x%X, want 0x%X", color, got, rooks)
		}
		if got := pos.Pieces[int(Queen)+int(color)]; got != queens {
			t.Errorf("Color %v Queens = 0x%X, want 0x%X", color, got, queens)
		}
		if got := pos.Pieces[int(King)+int(color)]; got != kings {
			t.Errorf("Color %v King = 0x%X, want 0x%X", color, got, kings)
		}
	}

	t.Run("Starting position", func(t *testing.T) {
		fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		pos, err := ParseFEN(fen)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		verifyPieces(t, pos, White,
			makeBitboard("a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2"), // Pawns
			makeBitboard("b1", "g1"),                                     // Knights
			makeBitboard("c1", "f1"),                                     // Bishops
			makeBitboard("a1", "h1"),                                     // Rooks
			makeBitboard("d1"),                                           // Queens
			makeBitboard("e1"),                                           // King
		)

		verifyPieces(t, pos, Black,
			makeBitboard("a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7"), // Pawns
			makeBitboard("b8", "g8"),                                     // Knights
			makeBitboard("c8", "f8"),                                     // Bishops
			makeBitboard("a8", "h8"),                                     // Rooks
			makeBitboard("d8"),                                           // Queens
			makeBitboard("e8"),                                           // King
		)

		if pos.SideToMove != White {
			t.Errorf("SideToMove = %v, want %v", pos.SideToMove, White)
		}
		expectedCastle := WhiteShort | WhiteLong | BlackShort | BlackLong
		if pos.CastlingRights != expectedCastle {
			t.Errorf("CastlingRights = %v, want %v", pos.CastlingRights, expectedCastle)
		}
		if pos.EnPassantSquare != NoSquare {
			t.Errorf("EnPassantSquare = %v, want %v", pos.EnPassantSquare, NoSquare)
		}
	})

	t.Run("Complex middlegame position", func(t *testing.T) {
		// A middlegame position:
		// - Black to move (b)
		// - Only White King-side castling allowed (K)
		// - En passant square is d3
		fen := "r2q1rk1/1pp1bppp/p1n1pn2/2P5/3Pp3/1NN5/PP2BPPP/R1BQK2R b K d3 0 11"
		pos, err := ParseFEN(fen)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		verifyPieces(t, pos, White,
			makeBitboard("a2", "b2", "f2", "g2", "h2", "c5", "d4"), // Pawns
			makeBitboard("b3", "c3"),                               // Knights
			makeBitboard("c1", "e2"),                               // Bishops
			makeBitboard("a1", "h1"),                               // Rooks
			makeBitboard("d1"),                                     // Queens
			makeBitboard("e1"),                                     // King
		)

		verifyPieces(t, pos, Black,
			makeBitboard("b7", "c7", "f7", "g7", "h7", "a6", "e6", "e4"), // Pawns
			makeBitboard("c6", "f6"),                                     // Knights
			makeBitboard("e7"),                                           // Bishops
			makeBitboard("a8", "f8"),                                     // Rooks
			makeBitboard("d8"),                                           // Queens
			makeBitboard("g8"),                                           // King
		)

		if pos.SideToMove != Black {
			t.Errorf("SideToMove = %v, want %v", pos.SideToMove, Black)
		}
		if pos.CastlingRights != WhiteShort {
			t.Errorf("CastlingRights = %v, want %v", pos.CastlingRights, WhiteShort)
		}
		expectedEP := GetNumBySquareName("d3")
		if pos.EnPassantSquare != expectedEP {
			t.Errorf("EnPassantSquare = %v, want %v", pos.EnPassantSquare, expectedEP)
		}
		if pos.HalfMoves != 0 || pos.CurrentTurn != 11 {
			t.Errorf("Moves = %d, %d, want 0, 11", pos.HalfMoves, pos.CurrentTurn)
		}
	})
}

func TestPosition_FEN_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		fen  string
	}{
		{
			name: "Starting position",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "Midgame position",
			fen:  "r1bqk2r/pp2bppp/2n1pn2/2p3B1/3pP3/3P1N2/PPPNBPPP/R2Q1RK1 b kq - 3 50",
		},
		{
			name: "Endgame position",
			fen:  "8/8/8/4k3/4P3/4K3/8/8 w - - 0 50",
		},
		{
			name: "En passant position",
			fen:  "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 50",
		},
		{
			name: "Alternating pieces and empty squares",
			fen:  "k7/1p1p1p1p/p1p1p1p1/1P1P1P1P/P1P1P1P1/8/8/7K w - - 0 50",
		},
		{
			name: "En passant on edge file",
			fen:  "k7/8/8/8/Pp6/8/8/7K b - a3 0 50",
		},
		{
			name: "En passant with no possible capturer",
			fen:  "k7/8/8/8/P6p/8/8/7K b - - 0 50",
		},
		{
			name: "Partial castling combinations",
			fen:  "r3k1r1/1p1n1ppp/p1pqbn2/3p4/3P4/2NBPN1P/PP3PP1/1R1QK2R w Kq - 2 13",
		},
		{
			name: "Multi-digit move counters",
			fen:  "8/8/8/4k3/4P3/4K3/8/8 w - - 99 150",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse setup FEN: %v", err)
			}

			got := pos.FEN()
			if got != tt.fen {
				t.Errorf("pos.FEN() = %q, want %q", got, tt.fen)
			}
		})
	}
}

func TestGetSquareNameByNum(t *testing.T) {
	tests := []struct {
		name string
		num  Square
		want string
	}{
		{name: "a8 square", num: 0, want: "a8"},
		{name: "h1 square", num: 63, want: "h1"},
		{name: "e4 square", num: 36, want: "e4"},
		{name: "NoSquare", num: NoSquare, want: "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSquareNameByNum(tt.num); got != tt.want {
				t.Errorf("GetSquareNameByNum(%d) = %q, want %q", tt.num, got, tt.want)
			}
		})
	}
}

func TestGetNumBySquareName(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want Square
	}{
		{name: "a8 square", str: "a8", want: 0},
		{name: "h1 square", str: "h1", want: 63},
		{name: "e4 square", str: "e4", want: 36},
		{name: "invalid length too short", str: "a", want: NoSquare},
		{name: "invalid length too long", str: "a12", want: NoSquare},
		{name: "wrong square format", str: "2f", want: NoSquare},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNumBySquareName(tt.str); got != tt.want {
				t.Errorf("GetNumBySquareName(%q) = %d, want %d", tt.str, got, tt.want)
			}
		})
	}
}

func TestPosition_FEN_Nil(t *testing.T) {
	var pos *Position
	if got := pos.FEN(); got != "" {
		t.Errorf("FEN() = %q, want %q", got, "")
	}
}

func TestPieceToFenChar_Invalid(t *testing.T) {
	char := pieceToFenChar(PieceType(99), White)
	if char != '?' {
		t.Errorf("pieceToFenChar() = %c, want %c", char, '?')
	}
}
