package engine

import (
	"fmt"
	"testing"
)

func TestBitboard_AddPiece(t *testing.T) {
	tests := []struct {
		name        string
		val         uint64
		squareNo    Square
		bitboardVal Bitboard
	}{
		{
			name:        "Adding piece to empty square",
			val:         0b_00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
			squareNo:    32,
			bitboardVal: 0b_00000000_00000000_00000000_00000001_00000000_00000000_11111111_00000000,
		},
		{
			name:        "Adding piece to non-empty square",
			val:         0b_00000000_00000000_00000000_00000001_00000000_00000000_11111111_00000000,
			squareNo:    32,
			bitboardVal: 0b_00000000_00000000_00000000_00000001_00000000_00000000_11111111_00000000,
		},
		{
			name:        "Adding piece to square 0 (a8)",
			val:         0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			squareNo:    0,
			bitboardVal: 0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
		},
		{
			name:        "Adding piece to square 63 (h1)",
			val:         0b_00000000_00000000_00000000_00000001_00000000_00000000_11111111_00000000,
			squareNo:    63,
			bitboardVal: 0b_10000000_00000000_00000000_00000001_00000000_00000000_11111111_00000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitboard := NewBitboard(tt.val)

			bitboard.AddPiece(tt.squareNo)

			if bitboard != tt.bitboardVal {
				t.Errorf("Bitboard(%064b).AddPiece(%v) = %064b, want %064b", tt.val, tt.squareNo, bitboard, tt.bitboardVal)
			}
		})
	}
}

func TestBitboard_String(t *testing.T) {
	tests := []struct {
		name   string
		val    uint64
		result string
	}{
		{
			name:   "Bitboard with zeros and ones",
			val:    0b_01000010_00000000_00000000_00111100_00000000_00000000_11111111_00000000,
			result: "0100001000000000000000000011110000000000000000001111111100000000",
		},
		{
			name:   "Empty bitboard",
			val:    0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			result: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:   "Full bitboard",
			val:    0b_11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
			result: "1111111111111111111111111111111111111111111111111111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitboard := NewBitboard(tt.val)

			bitboardStr := fmt.Sprint(bitboard)

			if bitboardStr != tt.result {
				t.Errorf("Bitboard(%064b).String() = %q, want %q", tt.val, bitboardStr, tt.result)
			}
		})
	}
}

func TestBitboard_Count(t *testing.T) {
	tests := []struct {
		name   string
		val    uint64
		result int
	}{
		{
			name:   "Bitboard with one rank as ones",
			val:    0b_11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			result: 8,
		},
		{
			name:   "Bitboard with one file as ones",
			val:    0b_10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000,
			result: 8,
		},
		{
			name:   "Bitboard with random squares included",
			val:    0b_01000010_01100000_00000000_00011000_00000000_00000000_00000000_10101010,
			result: 10,
		},
		{
			name:   "Empty bitboard",
			val:    0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			result: 0,
		},
		{
			name:   "Full bitboard",
			val:    0b_11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
			result: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitboard := NewBitboard(tt.val)

			count := bitboard.Count()

			if count != tt.result {
				t.Errorf("Bitboard(%064b).Count() = %d, want %d", tt.val, count, tt.result)
			}
		})
	}
}
