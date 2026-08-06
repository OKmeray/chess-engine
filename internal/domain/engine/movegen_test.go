package engine

import (
	"strings"
	"testing"
)

func TestGeneratePawnMoves(t *testing.T) {
	tests := []struct {
		name     string
		piece    Bitboard
		own      Bitboard
		enemy    Bitboard
		color    PieceColor
		epSquare Bitboard
		result   Bitboard
	}{
		// White Pawn Moves
		{
			name: "White pawn e3 single push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0, // 0 == empty board
			enemy: 0,
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e2 single and double push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e2 entirely blocked by own piece directly ahead",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			color:  White,
			result: 0,
		},
		{
			name: "White pawn e2 double push blocked by own, single push allowed",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e2 entirely blocked by enemy directly ahead",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color:  White,
			result: 0,
		},
		{
			name: "White pawn e2 double push blocked by enemy, single push allowed",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e4 both captures and single push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e5 single push and en passant capture",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			epSquare: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn a4 captures prevent H-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. X . . . . . X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X X . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn h4 captures prevent A-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . X .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e4 single push allowed but friendly capture prevented",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			color: White,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn e7 promotion",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			color: White,
			result: parseVisualBoard(`
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "White pawn d7 promotion and enemy capture, friendly capture prevented",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: parseVisualBoard(`
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: White,
			result: parseVisualBoard(`
				. . . X X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		// Black Pawn Moves
		{
			name: "Black pawn e6 single push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e7 single and double push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e7 entirely blocked by own piece directly ahead",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			color:  Black,
			result: 0,
		},
		{
			name: "Black pawn e7 double push blocked by own, single push allowed",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e7 entirely blocked by enemy directly ahead",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color:  Black,
			result: 0,
		},
		{
			name: "Black pawn e7 double push blocked by enemy, single push allowed",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e5 both captures and single push",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e4 single push and en passant capture",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: Black,
			epSquare: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . X . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X X . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn a5 captures prevent H-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. X . . . . . X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X X . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn h5 captures prevent A-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . X .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e5 single push allowed but friendly capture prevented",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Black pawn e2 promotion",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
			`),
		},
		{
			name: "Black pawn d2 forward promotion and enemy capture, friendly capture prevented",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . X . . . . .
			`),
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
			`),
			color: Black,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X X . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratePawnMoves(tt.piece, tt.own, tt.enemy, tt.color, tt.epSquare)
			if result != tt.result {
				t.Errorf("GeneratePawnMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		name   string
		piece  Bitboard
		own    Bitboard
		result Bitboard
	}{
		{
			name: "Knight e4 generates all 8 valid moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . X . . . X .
				. . . . . . . .
				. . X . . . X .
				. . . X . X . .
				. . . . . . . .
			`),
		},
		{
			name: "Knight a8 generates exactly 2 moves preventing edge wrapping",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . X . . . . .
				. X . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Knight h1 generates exactly 2 moves preventing edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X .
				. . . . . X . .
				. . . . . . . .
			`),
		},
		{
			name: "Knight a4 prevents G-file and H-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. X . . . . . .
				. . X . . . . .
				. . . . . . . .
				. . X . . . . .
				. X . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Knight b4 prevents H-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. X . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				X . X . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . X . . . .
				X . X . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Knight g4 prevents A-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . X . X
				. . . . X . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . X . X
				. . . . . . . .
			`),
		},
		{
			name: "Knight h4 prevents A-file and B-file edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . X .
				. . . . . X . .
				. . . . . . . .
				. . . . . X . .
				. . . . . . X .
				. . . . . . . .
			`),
		},
		{
			name: "Knight e4 generates 4 valid moves and prevents friendly captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X .
				. . . . . . . .
				. . X . . . X .
				. . . . . X . .
				. . . . . . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateKnightMoves(tt.piece, tt.own)
			if result != tt.result {
				t.Errorf("GenerateKnightMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	tests := []struct {
		name   string
		piece  Bitboard
		own    Bitboard
		enemy  Bitboard
		result Bitboard
	}{
		{
			name: "Bishop e5 on empty board generates full diagonals to edges",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. X . . . . . X
				. . X . . . X .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . X . . . X .
				. X . . . . . X
				X . . . . . . .
			`),
		},
		{
			name: "Bishop a8 on empty board generates single diagonal to opposite corner",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. X . . . . . .
				. . X . . . . .
				. . . X . . . .
				. . . . X . . .
				. . . . . X . .
				. . . . . . X .
				. . . . . . . X
			`),
		},
		{
			name: "Bishop a1 on empty board generates single diagonal to opposite corner",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . X
				. . . . . . X .
				. . . . . X . .
				. . . . X . . .
				. . . X . . . .
				. . X . . . . .
				. X . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Bishop e5 blocked by own pieces prevents captures and stops short",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . X
				. . . . . . X .
				. . . X . X . .
				. . . . . . . .
				. . . X . . . .
				. . X . . . . .
				. X . . . . . .
				X . . . . . . .
			`),
		},
		{
			name: "Bishop a8 blocked by own piece stops short",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. X . . . . . .
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Bishop e5 blocked by enemies includes captures but stops beyond",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . X
				. . X . . . X .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . X . . . . .
				. X . . . . . .
				X . . . . . . .
			`),
		},
		{
			name: "Bishop a8 blocked by enemy includes capture but stops beyond",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. X . . . . . .
				. . X . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Bishop e5 entirely enclosed by own pieces generates zero moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			result: 0,
		},
		{
			name: "Bishop e5 entirely enclosed by enemies generates exactly 4 captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Bishop e5 with multiple friendly blockers per ray correctly shadows pieces behind them",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. X . . . . . .
				. . X . . . X .
				. . . . . X . .
				. . . . . . . .
				. . . . . X . .
				. . X . . . . .
				. X . . . . . X
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Bishop e5 with multiple enemy blockers per ray includes first capture and shadows beyond",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. X . . . . . .
				. . X . . . X .
				. . . . . X . .
				. . . . . . . .
				. . . . . X . .
				. . X . . . . .
				. X . . . . . X
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . X . . . . .
				. . . X . X . .
				. . . . . . . .
				. . . X . X . .
				. . X . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBishopMoves(tt.piece, tt.own, tt.enemy)
			if result != tt.result {
				t.Errorf("GenerateBishopMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func TestGenerateRookMoves(t *testing.T) {
	tests := []struct {
		name   string
		piece  Bitboard
		own    Bitboard
		enemy  Bitboard
		result Bitboard
	}{
		{
			name: "Rook e5 on empty board generates full file and rank",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . X . . .
				. . . . X . . .
				. . . . X . . .
				X X X X . X X X
				. . . . X . . .
				. . . . X . . .
				. . . . X . . .
				. . . . X . . .
			`),
		},
		{
			name: "Rook a8 on empty board generates full file and rank to edges",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. X X X X X X X
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
			`),
		},
		{
			name: "Rook h1 on empty board generates full file and rank to edges",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . X
				. . . . . . . X
				. . . . . . . X
				. . . . . . . X
				. . . . . . . X
				. . . . . . . X
				. . . . . . . X
				X X X X X X X .
			`),
		},
		{
			name: "Rook e5 blocked by own pieces prevents captures and stops short",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
				. . X . . . X .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . . X X . .
				. . . X . . . .
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook a8 blocked by own piece prevents capture and stops short",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. X X . . . . .
				X . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook e5 blocked by enemies includes captures but stops beyond",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . X . . . X .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . X . . .
				. . X X . X X .
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook a8 blocked by enemy includes capture but stops beyond",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . X . . . .
				. . . . . . . .
				. . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. X X X . . . .
				X . . . . . . .
				X . . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook e5 entirely enclosed by own pieces generates zero moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . X . X . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			result: 0,
		},
		{
			name: "Rook e5 entirely enclosed by enemies generates exactly 4 captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . X . X . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . X . X . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook e5 with multiple friendly blockers per ray correctly shadows pieces behind them",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. X . X . . X X
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . X . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . X . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Rook e5 with multiple enemy blockers per ray includes first capture and shadows beyond",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. X . X . . X X
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . X . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . X . . .
				. . . . X . . .
				. . . X . X X .
				. . . . X . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateRookMoves(tt.piece, tt.own, tt.enemy)
			if result != tt.result {
				t.Errorf("GenerateRookMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	tests := []struct {
		name   string
		piece  Bitboard
		own    Bitboard
		enemy  Bitboard
		result Bitboard
	}{
		{
			name: "Queen e5 on empty board generates full rays in all 8 directions",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. X . . X . . X
				. . X . X . X .
				. . . X X X . .
				X X X X . X X X
				. . . X X X . .
				. . X . X . X .
				. X . . X . . X
				X . . . X . . .
			`),
		},
		{
			name: "Queen a8 on empty board generates full rays to edges",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. X X X X X X X
				X X . . . . . .
				X . X . . . . .
				X . . X . . . .
				X . . . X . . .
				X . . . . X . .
				X . . . . . X .
				X . . . . . . X
			`),
		},
		{
			name: "Queen e5 entirely enclosed by own pieces generates zero moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			result: 0,
		},
		{
			name: "Queen e5 entirely enclosed by enemies generates exactly 8 captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Queen a8 blocked by own pieces prevents captures and stops short",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . X . . . .
				. . . . . . . .
				. . X . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. X X . . . . .
				X X . . . . . .
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Queen e5 with multiple friendly blockers correctly shadows beyond them",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . X . X . X .
				. . . . . . . .
				. X . . . . X .
				. . . . . . . .
				. . X . X . . .
				. . . . . . . X
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . X X . X . .
				. . . X X X . .
				. . . . . . X .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "Queen e5 with multiple enemy blockers includes captures and shadows beyond",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . X . X . X .
				. . . . . . . .
				. X . . . . X .
				. . . . . . . .
				. . X . X . . .
				. . . . . . . X
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . X . X . X .
				. . . X X X . .
				. X X X . X X .
				. . . X X X . .
				. . X . X . X .
				. . . . . . . X
				. . . . . . . .
			`),
		},
		{
			name: "Queen e5 perfectly merges diagonal and orthogonal rays with enemy shadowing",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. X . . X . . X
				. . X . X . X .
				. . . . . . . .
				. X X . . X X .
				. . . . . . . .
				. . X . X . X .
				. X . . X . . X
				X . . . X . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . X . X . X .
				. . . X X X . .
				. . X X . X . .
				. . . X X X . .
				. . X . X . X .
				. . . . . . . .
				. . . . . . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateQueenMoves(tt.piece, tt.own, tt.enemy)
			if result != tt.result {
				t.Errorf("GenerateQueenMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func TestGenerateKingMoves(t *testing.T) {
	tests := []struct {
		name   string
		piece  Bitboard
		own    Bitboard
		enemy  Bitboard
		result Bitboard
	}{
		{
			name: "King e5 generates all 8 adjacent moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "King a8 generates exactly 3 moves preventing edge wrapping",
			piece: parseVisualBoard(`
				X . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. X . . . . . .
				X X . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "King h1 generates exactly 3 moves preventing edge wrapping",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . X
			`),
			own:   0,
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . X X
				. . . . . . X .
			`),
		},
		{
			name: "King e5 entirely enclosed by own pieces generates zero moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy:  0,
			result: 0,
		},
		{
			name: "King e5 entirely enclosed by enemies generates exactly 8 captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "King e5 blocked by 3 friendly pieces generates exactly 5 moves",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . X X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			enemy: 0,
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . . . X . .
				. . . . . X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
		{
			name: "King e5 adjacent to 3 enemies generates 8 moves including 3 captures",
			piece: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			own: 0,
			enemy: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
				. . . X . . . .
				. . . X X . . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
			result: parseVisualBoard(`
				. . . . . . . .
				. . . . . . . .
				. . . X X X . .
				. . . X . X . .
				. . . X X X . .
				. . . . . . . .
				. . . . . . . .
				. . . . . . . .
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateKingMoves(tt.piece, tt.own, tt.enemy)
			if result != tt.result {
				t.Errorf("GenerateKingMoves() = \n%v\nwant \n%v", result, tt.result)
			}
		})
	}
}

func parseVisualBoard(board string) Bitboard {
	var bb uint64
	clean := strings.ReplaceAll(board, " ", "")
	clean = strings.ReplaceAll(clean, "\t", "")
	clean = strings.ReplaceAll(clean, "\n", "")

	for i, r := range clean {
		if i >= 64 {
			break
		}
		if r == 'X' || r == 'x' {
			bitIndex := i
			bb |= (1 << bitIndex)
		}
	}
	return Bitboard(bb)
}
