package engine

import (
	"errors"
	"testing"
)

func TestNewMoveAndGetters(t *testing.T) {
	tests := []struct {
		name      string
		from      Square
		to        Square
		flags     uint16
		wantPromo PieceType
	}{
		{
			name:      "Quiet move",
			from:      52, // e2
			to:        36, // e4
			flags:     FlagDoublePawn,
			wantPromo: None,
		},
		{
			name:      "Capture move",
			from:      1,
			to:        18,
			flags:     FlagCapture,
			wantPromo: None,
		},
		{
			name:      "Promotion to Queen",
			from:      8,
			to:        0,
			flags:     FlagPromoQueen,
			wantPromo: Queen,
		},
		{
			name:      "Promotion to Knight with Capture",
			from:      15,
			to:        6,
			flags:     FlagPromoKnightCapture,
			wantPromo: Knight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMove(tt.from, tt.to, tt.flags)

			if got := m.From(); got != tt.from {
				t.Errorf("From() = %d, want %d", got, tt.from)
			}
			if got := m.To(); got != tt.to {
				t.Errorf("To() = %d, want %d", got, tt.to)
			}
			if got := m.Flags(); got != tt.flags {
				t.Errorf("Flags() = %d, want %d", got, tt.flags)
			}
			if got := m.Promotion(); got != tt.wantPromo {
				t.Errorf("Promotion() = %d, want %d", got, tt.wantPromo)
			}
		})
	}
}

func TestParseSAN(t *testing.T) {
	tests := []struct {
		name      string
		fen       string
		san       string
		wantFrom  Square
		wantTo    Square
		wantFlags uint16
		wantErr   error
	}{
		// --- Quiet Moves ---
		{
			name:      "White pawn one square move",
			fen:       "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			san:       "d3",
			wantFrom:  51, // d2
			wantTo:    43, // d3
			wantFlags: FlagQuiet,
		},
		{
			name:      "White pawn two square move",
			fen:       "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:       "e4",
			wantFrom:  52, // e2
			wantTo:    36, // e4
			wantFlags: FlagQuiet,
		},
		{
			name:      "Black pawn one square move",
			fen:       "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			san:       "d6",
			wantFrom:  11, // d7
			wantTo:    19, // d6
			wantFlags: FlagQuiet,
		},
		{
			name:      "Black pawn two square move",
			fen:       "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			san:       "e5",
			wantFrom:  12, // e7
			wantTo:    28, // e5
			wantFlags: FlagQuiet,
		},
		{
			name:      "White knight move",
			fen:       "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:       "Nf3",
			wantFrom:  62, // g1
			wantTo:    45, // f3
			wantFlags: FlagQuiet,
		},
		{
			name:      "Black bishop move",
			fen:       "rnbqkbnr/pppp1ppp/8/4p3/4P3/2N5/PPPP1PPP/R1BQKBNR b KQkq - 0 2",
			san:       "Bc5",
			wantFrom:  5,  // f8
			wantTo:    26, // c5
			wantFlags: FlagQuiet,
		},
		{
			name:      "White rook move",
			fen:       "r2qk2r/1p1n1pb1/p2p1np1/3Pp2p/8/1N2BP2/PPPQB1PP/R3K2R w KQkq - 4 13",
			san:       "Rc1",
			wantFrom:  56, // a1
			wantTo:    58, // c1
			wantFlags: FlagQuiet,
		},
		{
			name:      "Black queen move",
			fen:       "rnbqkbnr/pppp1ppp/8/4p3/4P3/2N5/PPPP1PPP/R1BQKBNR b KQkq - 0 2",
			san:       "Qf6",
			wantFrom:  3,  // d8
			wantTo:    21, // f6
			wantFlags: FlagQuiet,
		},
		{
			name:      "White king move",
			fen:       "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			san:       "Ke2",
			wantFrom:  60, // e1
			wantTo:    52, // e2
			wantFlags: FlagQuiet,
		},
		{
			name:      "White short castling",
			fen:       "r2qk2r/1p1n1pb1/p2p1np1/3Pp2p/8/1N2BP2/PPPQB1PP/R3K2R w KQkq - 2 13",
			san:       "O-O",
			wantFrom:  60, // e1
			wantTo:    62, // g1
			wantFlags: FlagKingCastle,
		},
		{
			name:      "White long castling",
			fen:       "r2qk2r/1p1n1pb1/p2p1np1/3Pp2p/8/1N2BP2/PPPQB1PP/R3K2R w KQkq - 2 13",
			san:       "O-O-O",
			wantFrom:  60, // e1
			wantTo:    58, // c1
			wantFlags: FlagQueenCastle,
		},
		{
			name:      "Black short castling",
			fen:       "r2qk2r/1p1n1pb1/p2p1np1/3Pp2p/8/1N2BP2/PPPQB1PP/2KR3R b kq - 3 13",
			san:       "O-O",
			wantFrom:  4, // e8
			wantTo:    6, // g8
			wantFlags: FlagKingCastle,
		},
		{
			name:      "Black long castling",
			fen:       "r3k2r/1pqn1pb1/p2p1np1/3Pp2p/8/1N2BP2/PPPQB1PP/1K1R3R b kq - 5 14",
			san:       "O-O-O",
			wantFrom:  4, // e8
			wantTo:    2, // c8
			wantFlags: FlagQueenCastle,
		},
		// --- Captures ---
		{
			name:      "White pawn captures black pawn",
			fen:       "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			san:       "exd5",
			wantFrom:  36, // e4
			wantTo:    27, // d5
			wantFlags: FlagCapture,
		},
		{
			name:      "Black pawn captures white pawn",
			fen:       "rnbqkbnr/ppp1pppp/8/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 2",
			san:       "dxe4",
			wantFrom:  27, // d5
			wantTo:    36, // e4
			wantFlags: FlagCapture,
		},
		{
			name:      "White knight captures black pawn",
			fen:       "rnbqkb1r/pppp1ppp/5n2/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3",
			san:       "Nxe5",
			wantFrom:  45, // f3
			wantTo:    28, // e5
			wantFlags: FlagCapture,
		},
		{
			name:      "Black bishop captures white knight",
			fen:       "rnbqk1nr/pppp1ppp/8/4p3/1b2P3/2N2N2/PPPP1PPP/R1BQKB1R b KQkq - 0 3",
			san:       "Bxc3",
			wantFrom:  33, // b4
			wantTo:    42, // c3
			wantFlags: FlagCapture,
		},
		{
			name:      "White queen captures black pawn",
			fen:       "rnbqkbnr/pppp1ppp/8/8/3p4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			san:       "Qxd4",
			wantFrom:  59, // d1
			wantTo:    35, // d4
			wantFlags: FlagCapture,
		},
		{
			name:      "White pawn captures en passant left",
			fen:       "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2",
			san:       "dxe6",
			wantFrom:  27, // d5
			wantTo:    20, // e6
			wantFlags: FlagCapture,
		},
		{
			name:      "White pawn captures en passant right",
			fen:       "rnbqkbnr/pppp1ppp/8/4pP2/8/8/PPPPP1PP/RNBQKBNR w KQkq e6 0 2",
			san:       "fxe6",
			wantFrom:  29, // f5
			wantTo:    20, // e6
			wantFlags: FlagCapture,
		},
		{
			name:      "Black pawn captures en passant left",
			fen:       "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 2",
			san:       "dxe3",
			wantFrom:  35, // d4
			wantTo:    44, // e3
			wantFlags: FlagCapture,
		},
		{
			name:      "Black pawn captures en passant right",
			fen:       "rnbqkbnr/ppppp1pp/8/8/4Pp2/8/PPPP2PP/RNBQKBNR b KQkq e3 0 2",
			san:       "fxe3",
			wantFrom:  37, // f4
			wantTo:    44, // e3
			wantFlags: FlagCapture,
		},
		// --- Disambiguation Moves ---
		{
			name:      "Disambiguation by file",
			fen:       "r1bqkb1r/pp2pppp/2n2n2/2pp4/3P1B2/4PN2/PPP2PPP/RN1QKB1R w KQkq - 1 5",
			san:       "Nbd2",
			wantFrom:  57, // b1
			wantTo:    51, // d2
			wantFlags: FlagQuiet,
		},
		{
			name:      "Disambiguation by file with capture",
			fen:       "r1bqk2r/ppp2ppp/2n5/8/1b1P4/5p2/PP1N1PPN/R1BQK2R w KQkq - 0 11",
			san:       "Nhxf3",
			wantFrom:  55, // h2
			wantTo:    45, // f3
			wantFlags: FlagCapture,
		},
		{
			name:      "Disambiguation by rank",
			fen:       "r5k1/1p2bppp/8/8/4p3/2B1R3/1P3PPP/r4RK1 b - - 0 25",
			san:       "R1a4",
			wantFrom:  56, // a1
			wantTo:    32, // a4
			wantFlags: FlagQuiet,
		},
		{
			name:      "Disambiguation by rank with capture",
			fen:       "r5k1/1p2bppp/P7/8/r3p3/2B1R3/1P3PPP/5RK1 b - - 0 25",
			san:       "R4xa6",
			wantFrom:  32, // a4
			wantTo:    16, // a6
			wantFlags: FlagCapture,
		},
		{
			name:      "Disambiguation by both rank and file",
			fen:       "1r1q1rk1/1p2pppp/8/8/8/Q7/5PPP/Q1Q3K1 w - - 0 41",
			san:       "Qa1b2",
			wantFrom:  56, // a1
			wantTo:    49, // b2
			wantFlags: FlagQuiet,
		},
		{
			name:      "Disambiguation by both with capture",
			fen:       "1r1q1rk1/1p2pppp/8/8/8/Q7/1p3PPP/Q1Q3K1 w - - 0 41",
			san:       "Qa1xb2",
			wantFrom:  56, // a1
			wantTo:    49, // b2
			wantFlags: FlagCapture,
		},
		// --- Pawn Promotions ---
		{
			name:      "White pawn promotion to queen",
			fen:       "8/4P2p/k4P2/8/8/p4K2/1p5P/8 w - - 0 45",
			san:       "e8=Q",
			wantFrom:  12, // e7
			wantTo:    4,  // e8
			wantFlags: FlagPromoQueen,
		},
		{
			name:      "Black pawn promotion to queen",
			fen:       "4Q3/7p/k4P2/8/8/p4K2/1p5P/8 b - - 0 45",
			san:       "b1=Q",
			wantFrom:  49, // b2
			wantTo:    57, // b1
			wantFlags: FlagPromoQueen,
		},
		{
			name:      "White pawn promotion to knight",
			fen:       "8/4P2p/k4P2/8/8/p4K2/1p5P/8 w - - 0 45",
			san:       "e8=N",
			wantFrom:  12, // e7
			wantTo:    4,  // e8
			wantFlags: FlagPromoKnight,
		},
		{
			name:      "Black pawn promotion to knight",
			fen:       "4Q3/7p/k4P2/8/8/p4K2/1p5P/8 b - - 0 45",
			san:       "b1=N",
			wantFrom:  49, // b2
			wantTo:    57, // b1
			wantFlags: FlagPromoKnight,
		},
		{
			name:      "White pawn promotion to rook",
			fen:       "8/4P2p/k4P2/8/8/p4K2/1p5P/8 w - - 0 45",
			san:       "e8=R",
			wantFrom:  12, // e7
			wantTo:    4,  // e8
			wantFlags: FlagPromoRook,
		},
		{
			name:      "Black pawn promotion to bishop",
			fen:       "4Q3/7p/k4P2/8/8/p4K2/1p5P/8 b - - 0 45",
			san:       "b1=B",
			wantFrom:  49, // b2
			wantTo:    57, // b1
			wantFlags: FlagPromoBishop,
		},
		{
			name:      "White capture with promotion to queen",
			fen:       "3n4/4P2p/k4P2/8/8/p4K2/1p5P/2N5 w - - 0 45",
			san:       "exd8=Q",
			wantFrom:  12, // f7
			wantTo:    3,  // g8
			wantFlags: FlagPromoQueenCapture,
		},
		{
			name:      "Black capture with promotion to knight",
			fen:       "3Q4/7p/k4P2/8/8/p4K2/1p5P/2N5 b - - 0 45",
			san:       "bxc1=N",
			wantFrom:  49, // b2
			wantTo:    58, // c1
			wantFlags: FlagPromoKnightCapture,
		},
		// --- Check and Checkmate Modifiers ---
		{
			name:      "Black move with check",
			fen:       "rnbqkb1r/ppp2ppp/4pn2/3p4/2PP4/5NP1/PP2PP1P/RNBQKB1R b KQkq - 1 3",
			san:       "Bb4+",
			wantFrom:  5,  // f8
			wantTo:    33, // b4
			wantFlags: FlagQuiet,
		},
		{
			name:      "White move with checkmate",
			fen:       "rnbqkbnr/ppppp2p/5p2/6p1/4P3/2N5/PPPP1PPP/R1BQKBNR w KQkq - 0 3",
			san:       "Qh5#",
			wantFrom:  59, // d1
			wantTo:    31, // h5
			wantFlags: FlagQuiet,
		},
		{
			name:      "Capture with check",
			fen:       "rnbqk2r/ppp2ppp/4pn2/3p4/1bPP4/2N2NP1/PP2PP1P/R1BQKB1R b KQkq - 3 4",
			san:       "Bxc3+",
			wantFrom:  33, // b4
			wantTo:    42, // c3
			wantFlags: FlagCapture,
		},
		{
			name:      "Capture with checkmate",
			fen:       "5rrk/1R5p/8/5B2/8/8/P4PPP/6K1 w - - 0 32",
			san:       "Rxh7#",
			wantFrom:  9,  // b7
			wantTo:    15, // h7
			wantFlags: FlagCapture,
		},
		{
			name:      "Disambiguation with check",
			fen:       "8/pr3p1p/4p3/8/2k5/5N2/P4PPP/1N2K3 w - - 0 20",
			san:       "Nbd2+",
			wantFrom:  57, // b1
			wantTo:    51, // d2
			wantFlags: FlagQuiet,
		},
		{
			name:      "Disambiguation with checkmate",
			fen:       "1r6/p4p1p/2B5/2p5/2k5/P4N2/2K2PPP/1N6 w - - 0 20",
			san:       "Nbd2+",
			wantFrom:  57, // b1
			wantTo:    51, // d2
			wantFlags: FlagQuiet,
		},
		{
			name:      "Short castling with check",
			fen:       "6r1/8/8/p4kp1/P6p/7P/6P1/4K2R w K - 0 51",
			san:       "O-O+",
			wantFrom:  60, // e1
			wantTo:    62, // g1
			wantFlags: FlagKingCastle,
		},
		{
			name:      "Long castling with checkmate",
			fen:       "1r4r1/8/n7/p2k1Pp1/P6p/2R2P1P/1PP4B/R3K3 w Q - 1 51",
			san:       "O-O-O#",
			wantFrom:  60, // e1
			wantTo:    58, // c1
			wantFlags: FlagQueenCastle,
		},
		{
			name:      "Promotion with check",
			fen:       "k7/4P2p/5P2/8/8/p4K2/1p5P/8 w - - 0 45",
			san:       "e8=Q+",
			wantFrom:  12, // e7
			wantTo:    4,  // e8
			wantFlags: FlagPromoQueen,
		},
		{
			name:      "Promotion with checkmate",
			fen:       "k7/pp2P2p/5P2/8/8/5K2/7P/8 w - - 0 45",
			san:       "e8=R#",
			wantFrom:  12, // e7
			wantTo:    4,  // e8
			wantFlags: FlagPromoRook,
		},
		{
			name:      "Promotion via capture with check",
			fen:       "k4n2/4P2p/5P2/8/8/p4K2/1p5P/8 w - - 0 45",
			san:       "exf8=Q+",
			wantFrom:  12, // e7
			wantTo:    5,  // f8
			wantFlags: FlagPromoQueenCapture,
		},
		{
			name:      "Promotion via capture with checkmate",
			fen:       "k4n2/pp2P2p/5P2/8/8/5K2/7P/8 w - - 0 45",
			san:       "exf8=Q#",
			wantFrom:  12, // e7
			wantTo:    5,  // f8
			wantFlags: FlagPromoQueenCapture,
		},

		// --- Error Paths (Sad Paths) ---
		{
			name:    "Error san string too short",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "e",
			wantErr: ErrSANTooShort,
		},
		{
			name:    "Error invalid san move format (+)",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "+",
			wantErr: ErrSANTooShort,
		},
		{
			name:    "Error invalid san move format after stripping plus",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "e+",
			wantErr: ErrInvalidSANFormat,
		},
		{
			name:    "Error invalid promotion piece",
			fen:     "8/4P3/3k4/8/8/8/8/4K3 w - - 0 1",
			san:     "e8=K",
			wantErr: ErrInvalidPromotion,
		},
		{
			name:    "Error invalid promotion piece character in san",
			fen:     "8/4P3/3k4/8/8/8/8/4K3 w - - 0 1",
			san:     "e8=Z",
			wantErr: ErrInvalidPromotion,
		},
		{
			name:      "Promotion without equals sign",
			fen:       "8/4P3/3k4/8/8/8/8/4K3 w - - 0 1",
			san:       "e8Q",
			wantFrom:  12, // e7
			wantTo:    4,  // e8  // e8
			wantFlags: FlagPromoQueen,
		},
		{
			name:    "Error string too short after stripping promotion",
			fen:     "8/4P3/8/6k1/8/8/8/4K3 w - - 0 1",
			san:     "e=Q",
			wantErr: ErrSANTooShort,
		},
		{
			name:    "Error invalid destination square",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "e9",
			wantErr: ErrInvalidDestination,
		},
		{
			name:    "Error invalid piece character",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "Zf3",
			wantErr: ErrInvalidDisambiguation,
		},
		{
			name:    "Error invalid disambiguation format",
			fen:     "r1bqkb1r/pp2pppp/2n2n2/2pp4/3P1B2/4PN2/PPP2PPP/RN1QKB1R w KQkq - 1 5",
			san:     "Nzd2",
			wantErr: ErrInvalidDisambiguation,
		},
		{
			name:    "Error invalid double disambiguation",
			fen:     "4k3/8/8/8/8/Q7/8/Q1Q1K3 w - - 0 1",
			san:     "Qv1b2",
			wantErr: ErrInvalidDoubleDisambiguation,
		},
		{
			name:    "Error invalid disambiguation length",
			fen:     "4k3/8/8/8/8/Q7/8/Q1Q1K3 w - - 0 1",
			san:     "Qa12b2",
			wantErr: ErrInvalidDisambiguation,
		},
		{
			name:    "Error impossible pawn move",
			fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			san:     "e5", // Pawn cannot move 3 squares
			wantErr: ErrIllegalOrAmbiguousMove,
		},
		{
			name:    "Error piece not found on board",
			fen:     "4k3/q7/8/8/8/8/6P1/4K3 w - - 0 1", // No white queen
			san:     "Qa3",
			wantErr: ErrIllegalOrAmbiguousMove,
		},
		{
			name:    "Error piece on wrong file",
			fen:     "4k3/8/8/8/R7/8/8/4K3 w - - 0 1", // Rook on a4, asking for b file
			san:     "Rbb3",
			wantErr: ErrIllegalOrAmbiguousMove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			move, err := ParseSAN(pos, tt.san)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseSAN(%q) error = %v, wantErr %v", tt.san, err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if got := move.From(); got != tt.wantFrom {
					t.Errorf("From() = %v, want %v", got, tt.wantFrom)
				}
				if got := move.To(); got != tt.wantTo {
					t.Errorf("To() = %v, want %v", got, tt.wantTo)
				}
				if got := move.Flags(); got != tt.wantFlags {
					t.Errorf("Flags() = %v, want %v", got, tt.wantFlags)
				}
			}
		})
	}
}

func TestParseSANNilPosition(t *testing.T) {
	_, err := ParseSAN(nil, "e4")
	if !errors.Is(err, ErrNilPosition) {
		t.Errorf("ParseSAN(nil, %q) error = %v, wantErr %v", "e4", err, ErrNilPosition)
	}
}
