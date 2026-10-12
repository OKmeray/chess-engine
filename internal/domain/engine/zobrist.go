package engine

import (
	"math/rand"
)

var (
	// ZobristTable holds random numbers for [color][piece][square]
	ZobristTable [2][6][64]uint64
	// ZobristSideToMove is the random number to XOR when it's Black's turn to move.
	ZobristSideToMove uint64
	// ZobristCastling holds random numbers for all 16 possible castling right combinations.
	ZobristCastling [16]uint64
	// ZobristEnPassant holds random numbers for the 8 possible files an en passant target could be on.
	ZobristEnPassant [8]uint64
)

func init() {
	InitZobrist()
}

// InitZobrist populates the Zobrist tables with deterministic random numbers.
func InitZobrist() {
	// A fixed seed for strict reproducibility.
	rng := rand.New(rand.NewSource(21212121))

	// Pieces
	for c := range 2 {
		for p := range 6 {
			for sq := range 64 {
				ZobristTable[c][p][sq] = rng.Uint64()
			}
		}
	}

	// Side to move
	ZobristSideToMove = rng.Uint64()

	// Castling rights (4 bits = 16 combinations)
	for i := range 16 {
		ZobristCastling[i] = rng.Uint64()
	}

	// En Passant files (8 files)
	for i := range 8 {
		ZobristEnPassant[i] = rng.Uint64()
	}
}

// ComputeHash calculates the Zobrist hash of the position from scratch.
// This is relatively slow and should only be used when parsing a FEN.
func (p *Position) ComputeHash() uint64 {
	var hash uint64

	for sq := range Square(64) {
		pt, color := p.GetPieceAndColorBySquare(sq)
		if pt != None {
			cIdx := color / 6 // Convert White(6)->1, Black(0)->0
			hash ^= ZobristTable[cIdx][pt][sq]
		}
	}

	if p.SideToMove == Black {
		hash ^= ZobristSideToMove
	}

	// CastlingRights is a bitmask (0-15), which maps directly to the ZobristCastling array index.
	hash ^= ZobristCastling[int(p.CastlingRights)]

	if p.EnPassantSquare != NoSquare {
		file := int(p.EnPassantSquare % 8)
		hash ^= ZobristEnPassant[file]
	}

	return hash
}
