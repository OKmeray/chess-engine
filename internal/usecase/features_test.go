package usecase

import (
	"testing"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// assertGlobalFeature ensures a global feature (like castling) is identical on all 64 squares.
func assertGlobalFeature(t *testing.T, out []float32, feature int, want float32, label string) {
	t.Helper()

	for sq := 0; sq < 64; sq++ {
		if got := out[sq*20+feature]; got != want {
			t.Errorf("ExtractTransformerFeatures() %s out[sq=%d, feature=%d] = %f, want %f",
				label, sq, feature, got, want)
		}
	}
}

// assertOnlyOneHot verifies exactly the specified squares are 1.0, and all others are strictly 0.0.
func assertOnlyOneHot(t *testing.T, out []float32, feature int, onSquares map[int]bool, label string) {
	t.Helper()

	for sq := 0; sq < 64; sq++ {
		want := float32(0.0)
		if onSquares[sq] {
			want = 1.0
		}
		if got := out[sq*20+feature]; got != want {
			t.Errorf("ExtractTransformerFeatures() %s out[sq=%d, feature=%d] = %f, want %f",
				label, sq, feature, got, want)
		}
	}
}

// extract parses a FEN and runs ExtractTransformerFeatures on a garbage-filled
// array to ensure all 1280 slots are safely initialized and overwritten.
func extract(t *testing.T, fen string) []float32 {
	t.Helper()

	pos, err := engine.ParseFEN(fen)
	if err != nil {
		t.Fatalf("failed to parse FEN %q: %v", fen, err)
	}

	out := make([]float32, 1280)
	// Pre-fill with garbage to ensure the extractor fully overwrites every slot.
	for i := range out {
		out[i] = -9.0
	}

	ExtractTransformerFeatures(pos, out)

	return out
}

func TestExtractTransformerFeatures_StartingPosition_Pieces(t *testing.T) {
	outWhite := extract(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	outBlack := extract(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")

	expectedPieces := make([]float32, 64*12)

	// index convention: 0=BP 1=BN 2=BB 3=BR 4=BQ 5=BK 6=WP 7=WN 8=WB 9=WR 10=WQ 11=WK
	set := func(sq int, piece int) {
		expectedPieces[sq*12+piece] = 1.0
	}

	// Black back rank (sq 0-7)
	set(0, 3) // a8 Black Rook
	set(1, 1) // b8 Black Knight
	set(2, 2) // c8 Black Bishop
	set(3, 4) // d8 Black Queen
	set(4, 5) // e8 Black King
	set(5, 2) // f8 Black Bishop
	set(6, 1) // g8 Black Knight
	set(7, 3) // h8 Black Rook

	// Black pawns (sq 8-15)
	for sq := 8; sq <= 15; sq++ {
		set(sq, 0)
	}

	// White pawns (sq 48-55)
	for sq := 48; sq <= 55; sq++ {
		set(sq, 6)
	}

	// White back rank (sq 56-63)
	set(56, 9)  // a1 White Rook
	set(57, 7)  // b1 White Knight
	set(58, 8)  // c1 White Bishop
	set(59, 10) // d1 White Queen
	set(60, 11) // e1 White King
	set(61, 8)  // f1 White Bishop
	set(62, 7)  // g1 White Knight
	set(63, 9)  // h1 White Rook

	for sq := 0; sq < 64; sq++ {
		for f := 0; f < 12; f++ {
			expectedVal := expectedPieces[sq*12+f]

			gotWhite := outWhite[sq*20+f]
			if gotWhite != expectedVal {
				t.Errorf("Unmirrored (White to move): out[sq=%d, feature=%d] = %f, want %f", sq, f, gotWhite, expectedVal)
			}

			gotBlack := outBlack[sq*20+f]
			if gotBlack != expectedVal {
				t.Errorf("Mirrored (Black to move): out[sq=%d, feature=%d] = %f, want %f", sq, f, gotBlack, expectedVal)
			}
		}
	}
}

func TestExtractTransformerFeatures_Castling_NonMirrored(t *testing.T) {
	out := extract(t, "r3k3/6p1/7p/8/8/P7/1P6/4K2R w Kq - 5 51")
	assertGlobalFeature(t, out, 14, 1.0, "wsc")
	assertGlobalFeature(t, out, 15, 0.0, "wlc")
	assertGlobalFeature(t, out, 16, 0.0, "bsc")
	assertGlobalFeature(t, out, 17, 1.0, "blc")
}

func TestExtractTransformerFeatures_Castling_Mirrored(t *testing.T) {
	out := extract(t, "r6r/3k2p1/7p/8/8/P7/1P6/R3K2R b KQ - 5 51")
	assertGlobalFeature(t, out, 14, 0.0, "wsc")
	assertGlobalFeature(t, out, 15, 0.0, "wlc")
	assertGlobalFeature(t, out, 16, 1.0, "bsc")
	assertGlobalFeature(t, out, 17, 1.0, "blc")
}

func TestExtractTransformerFeatures_Castling_None(t *testing.T) {
	out := extract(t, "r6r/3k2p1/7p/8/8/P7/1P6/R3K2R w - - 5 51")
	for _, f := range []int{14, 15, 16, 17} {
		assertGlobalFeature(t, out, f, 0.0, "no castling rights")
	}
}

func TestExtractTransformerFeatures_EnPassant_NonMirrored(t *testing.T) {
	out := extract(t, "8/5p1k/8/3Pp3/8/6P1/3K4/8 w - e6 0 52")
	assertOnlyOneHot(t, out, 18, map[int]bool{20: true}, "en passant (non-mirrored)")
}

func TestExtractTransformerFeatures_EnPassant_Mirrored(t *testing.T) {
	out := extract(t, "8/3k1p2/8/8/2pP4/8/4P3/2K5 b - d3 0 52")
	assertOnlyOneHot(t, out, 18, map[int]bool{19: true}, "en passant (mirrored)") // 43 mirrored to 19
}

func TestExtractTransformerFeatures_EnPassant_None(t *testing.T) {
	out := extract(t, "8/3kb3/4pp2/1Pp2n2/2P5/8/1R1KP3/8 w - - 1 52")
	assertOnlyOneHot(t, out, 18, map[int]bool{}, "no en passant")
}

func TestExtractTransformerFeatures_Halfmoves_FullBitPattern(t *testing.T) {
	out := extract(t, "4k3/8/8/3p4/8/4P3/8/4K3 w - - 5 61") // 5 in binary is ...00000101
	expectedOnes := map[int]bool{61: true, 63: true}        // inclusive bits in number 5
	for sq := 0; sq < 64; sq++ {
		want := float32(0.0)
		if expectedOnes[sq] {
			want = 1.0
		}
		if got := out[sq*20+19]; got != want {
			t.Errorf("ExtractTransformerFeatures() halfmoves out[sq=%d] = %f, want %f", sq, got, want)
		}
	}
}

func TestExtractTransformerFeatures_Halfmoves_Zero(t *testing.T) {
	out := extract(t, "4k3/8/8/3p4/8/4P3/8/4K3 w - - 0 61")
	for sq := 0; sq < 64; sq++ {
		if got := out[sq*20+19]; got != 0.0 {
			t.Errorf("ExtractTransformerFeatures() halfmoves out[sq=%d] = %f, want 0.0", sq, got)
		}
	}
}

func TestExtractTransformerFeatures_RepetitionFeaturesAreZero(t *testing.T) {
	// Features 12-13 (repetition) are TODO.
	// Ensure they output 0.0 to prevent buffer garbage leakage.
	out := extract(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	for sq := 0; sq < 64; sq++ {
		for _, f := range []int{12, 13} {
			if got := out[sq*20+f]; got != 0.0 {
				t.Errorf("ExtractTransformerFeatures() repetition out[sq=%d, feature=%d] = %f, want 0.0", sq, f, got)
			}
		}
	}
}
