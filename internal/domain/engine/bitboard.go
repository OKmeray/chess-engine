package engine

import (
	"fmt"
	"math/bits"
)

// Square represents a position on the chessboard (0-63)
type Square int

const (
	// NoSquare is a sentinel value indicating the absence of a valid square
	NoSquare Square = 64
)

// Bitboard represents board property as 64-bit number
type Bitboard uint64

// NewBitboard returns new bitboard
func NewBitboard(value uint64) Bitboard {
	return Bitboard(value)
}

// AddPiece adds piece to the bitboard
func (b *Bitboard) AddPiece(square Square) {
	*b |= 1 << uint64(square)
}

// String return bitboard as a string of 64 bit symbols
func (b Bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}

// Count returns the number of set bits on the bitboard
func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}
