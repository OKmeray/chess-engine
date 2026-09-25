package usecase

import (
	"testing"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

func TestEncodeMove(t *testing.T) {
	tests := []struct {
		name   string
		move   engine.Move
		mirror bool
		want   int
	}{
		// ====================
		// Slide Moves
		// ====================
		{
			name:   "Slide N: e2-e4", // Path: 12 -> 28 | Action: Slide N (2 steps) | Plane: 1
			move:   engine.NewMove(12, 28, engine.FlagQuiet),
			mirror: false,
			want:   12*76 + 1,
		},
		{
			name:   "Slide S: a8-a1", // Path: 56 -> 0 | Action: Slide S (7 steps) | Plane: 28 (base) + 6 (offset) = 34
			move:   engine.NewMove(56, 0, engine.FlagQuiet),
			mirror: false,
			want:   56*76 + 34,
		},
		{
			name:   "Slide E: a1-h1", // Path: 0 -> 7 | Action: Slide E (7 steps) | Plane: 14 (base) + 6 (offset) = 20
			move:   engine.NewMove(0, 7, engine.FlagQuiet),
			mirror: false,
			want:   0*76 + 20,
		},
		{
			name:   "Slide W: h1-a1", // Path: 7 -> 0 | Action: Slide W (7 steps) | Plane: 42 (base) + 6 (offset) = 48
			move:   engine.NewMove(7, 0, engine.FlagQuiet),
			mirror: false,
			want:   7*76 + 48,
		},
		{
			name:   "Slide NE: a1-h8", // Path: 0 -> 63 | Action: Slide NE (7 steps) | Plane: 7 (base) + 6 (offset) = 13
			move:   engine.NewMove(0, 63, engine.FlagQuiet),
			mirror: false,
			want:   0*76 + 13,
		},
		{
			name:   "Slide NW: h1-a8", // Path: 7 -> 56 | Action: Slide NW (7 steps) | Plane: 49 (base) + 6 (offset) = 55
			move:   engine.NewMove(7, 56, engine.FlagQuiet),
			mirror: false,
			want:   7*76 + 55,
		},
		{
			name:   "Slide SE: a8-h1", // Path: 56 -> 7 | Action: Slide SE (7 steps) | Plane: 21 (base) + 6 (offset) = 27
			move:   engine.NewMove(56, 7, engine.FlagQuiet),
			mirror: false,
			want:   56*76 + 27,
		},
		{
			name:   "Slide SW: h8-a1", // Path: 63 -> 0 | Action: Slide SW (7 steps) | Plane: 35 (base) + 6 (offset) = 41
			move:   engine.NewMove(63, 0, engine.FlagQuiet),
			mirror: false,
			want:   63*76 + 41,
		},

		// ====================
		// Knight Moves
		// ====================
		{
			name:   "Knight NE: b1-c3", // Path: 1 -> 18 | Action: Knight N-E | Plane: 56
			move:   engine.NewMove(1, 18, engine.FlagQuiet),
			mirror: false,
			want:   1*76 + 56,
		},
		{
			name:   "Knight NW: e2-d4", // Path: 12 -> 27 | Action: Knight N-W | Plane: 57
			move:   engine.NewMove(12, 27, engine.FlagQuiet),
			mirror: false,
			want:   12*76 + 57,
		},
		{
			name:   "Knight SE: e4-f2", // Path: 28 -> 13 | Action: Knight S-E | Plane: 58
			move:   engine.NewMove(28, 13, engine.FlagQuiet),
			mirror: false,
			want:   28*76 + 58,
		},
		{
			name:   "Knight SW: e4-d2", // Path: 28 -> 11 | Action: Knight S-W | Plane: 59
			move:   engine.NewMove(28, 11, engine.FlagQuiet),
			mirror: false,
			want:   28*76 + 59,
		},
		{
			name:   "Knight EN: e2-g3", // Path: 12 -> 22 | Action: Knight E-N | Plane: 60
			move:   engine.NewMove(12, 22, engine.FlagQuiet),
			mirror: false,
			want:   12*76 + 60,
		},
		{
			name:   "Knight ES: e4-g3", // Path: 28 -> 22 | Action: Knight E-S | Plane: 61
			move:   engine.NewMove(28, 22, engine.FlagQuiet),
			mirror: false,
			want:   28*76 + 61,
		},
		{
			name:   "Knight WN: e2-c3", // Path: 12 -> 18 | Action: Knight W-N | Plane: 62
			move:   engine.NewMove(12, 18, engine.FlagQuiet),
			mirror: false,
			want:   12*76 + 62,
		},
		{
			name:   "Knight WS: e4-c3", // Path: 28 -> 18 | Action: Knight W-S | Plane: 63
			move:   engine.NewMove(28, 18, engine.FlagQuiet),
			mirror: false,
			want:   28*76 + 63,
		},

		// ====================
		// Promotions
		// ====================
		{
			name:   "Promo Queen: e7-e8=Q", // Path: 52 -> 60 | Action: Promo North to Queen | Plane: 68
			move:   engine.NewMove(52, 60, engine.FlagPromoQueen),
			mirror: false,
			want:   52*76 + 68,
		},
		{
			name:   "Promo Queen Capture: e7-f8=Q", // Path: 52 -> 61 | Action: Promo Capture N-E to Queen | Plane: 72
			move:   engine.NewMove(52, 61, engine.FlagPromoQueenCapture),
			mirror: false,
			want:   52*76 + 72,
		},
		{
			name:   "Promo Rook: e7-e8=R", // Path: 52 -> 60 | Action: Promo North to Rook | Plane: 69
			move:   engine.NewMove(52, 60, engine.FlagPromoRook),
			mirror: false,
			want:   52*76 + 69,
		},
		{
			name:   "Promo Rook Capture: e7-d8=R", // Path: 52 -> 59 | Action: Promo Capture N-W to Rook | Plane: 65
			move:   engine.NewMove(52, 59, engine.FlagPromoRookCapture),
			mirror: false,
			want:   52*76 + 65,
		},
		{
			name:   "Promo Bishop: e7-e8=B", // Path: 52 -> 60 | Action: Promo North to Bishop | Plane: 70
			move:   engine.NewMove(52, 60, engine.FlagPromoBishop),
			mirror: false,
			want:   52*76 + 70,
		},
		{
			name:   "Promo Bishop Capture: e7-f8=B", // Path: 52 -> 61 | Action: Promo Capture N-E to Bishop | Plane: 74
			move:   engine.NewMove(52, 61, engine.FlagPromoBishopCapture),
			mirror: false,
			want:   52*76 + 74,
		},
		{
			name:   "Promo Knight: e7-e8=N", // Path: 52 -> 60 | Action: Promo North to Knight | Plane: 71
			move:   engine.NewMove(52, 60, engine.FlagPromoKnight),
			mirror: false,
			want:   52*76 + 71,
		},
		{
			name:   "Promo Knight Capture: e7-d8=N", // Path: 52 -> 59 | Action: Promo Capture N-W to Knight | Plane: 67
			move:   engine.NewMove(52, 59, engine.FlagPromoKnightCapture),
			mirror: false,
			want:   52*76 + 67,
		},

		// ====================
		// Mirroring
		// ====================
		{
			name:   "Mirrored Slide: h8-h1", // Path: 63 -> 7 | Action: Slide North (7 steps) mirrored | Plane: 6
			move:   engine.NewMove(63, 7, engine.FlagQuiet),
			mirror: true,
			want:   7*76 + 6,
		},
		{
			name:   "Mirrored Promo: e2-e1=Q", // Path: 12 -> 4 | Action: Mirrors to e7(52) -> e8(60). Promo North to Queen | Plane: 68
			move:   engine.NewMove(12, 4, engine.FlagPromoQueen),
			mirror: true,
			want:   52*76 + 68,
		},
		{
			name:   "Black Promotion e2-e1=Q Unmirrored", // Path: 12 -> 4 | Action: Promo North to Queen but unmirrored -> S Q | Plane: 0
			move:   engine.NewMove(12, 4, engine.FlagPromoQueen),
			mirror: false,
			want:   12*76 + 0, // Not initialized in map, falls back to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeMove(tt.move, tt.mirror)
			if got != tt.want {
				t.Errorf("EncodeMove() = %v, want %v", got, tt.want)
			}
		})
	}
}
