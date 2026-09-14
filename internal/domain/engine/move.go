package engine

import (
	"errors"
	"fmt"
	"math/bits"
)

var (
	ErrNilPosition                 = errors.New("position cannot be nil")
	ErrSANTooShort                 = errors.New("SAN string too short")
	ErrInvalidSANFormat            = errors.New("invalid SAN move format")
	ErrInvalidPromotion            = errors.New("invalid promotion piece in SAN")
	ErrInvalidDestination          = errors.New("invalid destination square in SAN")
	ErrInvalidPiece                = errors.New("invalid piece character in SAN")
	ErrInvalidDisambiguation       = errors.New("invalid disambiguation format in SAN")
	ErrInvalidDoubleDisambiguation = errors.New("invalid double disambiguation in SAN")
	ErrIllegalOrAmbiguousMove      = errors.New("illegal move or ambiguous SAN")
)

// Move represents a chess move encoded into a 16-bit integer.
//
// Bit Layout:
// - Bits 0-5 (6 bits)  : Destination Square (0-63)
// - Bits 6-11 (6 bits) : Origin Square (0-63)
// - Bits 12-15 (4 bits): Move Flags (Promotion, Castling, En Passant, etc.)
type Move uint16

// Move flags
const (
	FlagQuiet       uint16 = 0b0000
	FlagDoublePawn  uint16 = 0b0001
	FlagKingCastle  uint16 = 0b0010
	FlagQueenCastle uint16 = 0b0011
	FlagCapture     uint16 = 0b0100
	FlagEnPassant   uint16 = 0b0101

	// Promotion flags
	FlagPromoKnight uint16 = 0b1000
	FlagPromoBishop uint16 = 0b1001
	FlagPromoRook   uint16 = 0b1010
	FlagPromoQueen  uint16 = 0b1011

	// Capture Promotion flags
	FlagPromoKnightCapture uint16 = 0b1100 // the same as FlagCapture + FlagPromoKnight
	FlagPromoBishopCapture uint16 = 0b1101 // the same as FlagCapture + FlagPromoBishop
	FlagPromoRookCapture   uint16 = 0b1110 // the same as FlagCapture + FlagPromoRook
	FlagPromoQueenCapture  uint16 = 0b1111 // the same as FlagCapture + FlagPromoQueen
)

const (
	squareMask       = 0x3F // Isolates the lowest 6 bits
	fromSquareOffset = 6
	flagOffset       = 12
)

// NewMove creates a new bit-packed move.
func NewMove(from, to Square, flags uint16) Move {
	return Move(to) | (Move(from) << fromSquareOffset) | (Move(flags) << flagOffset)
}

// From returns the origin square (0-63).
func (m Move) From() Square {
	return Square((m >> fromSquareOffset) & squareMask)
}

// To returns the destination square (0-63).
func (m Move) To() Square {
	return Square(m & squareMask)
}

// Flags returns the 4-bit move flag.
func (m Move) Flags() uint16 {
	return uint16(m >> flagOffset)
}

// Promotion returns the piece type to promote to, or None if not a promotion.
func (m Move) Promotion() PieceType {
	flags := m.Flags()
	if flags >= FlagPromoKnight {
		// Maps promotion flags 8-15 to PieceType values 1-4 (Knight-Queen).
		return PieceType((flags & 0b0011) + 1)
	}
	return None
}

// ParseSAN converts a Standard Algebraic Notation (SAN) string into a Move.
func ParseSAN(pos *Position, san string) (Move, error) {
	if pos == nil {
		return Move(0), ErrNilPosition
	}
	if len(san) < 2 {
		return Move(0), ErrSANTooShort
	}

	color := pos.SideToMove
	originalSan := san

	lastByte := san[len(san)-1]
	if lastByte == '+' || lastByte == '#' {
		san = san[:len(san)-1]
	}

	if san == "O-O" {
		if color == White {
			return NewMove(60, 62, FlagKingCastle), nil
		}
		return NewMove(4, 6, FlagKingCastle), nil
	}
	if san == "O-O-O" {
		if color == White {
			return NewMove(60, 58, FlagQueenCastle), nil
		}
		return NewMove(4, 2, FlagQueenCastle), nil
	}

	if len(san) < 2 {
		return Move(0), fmt.Errorf("%w: %s", ErrInvalidSANFormat, originalSan)
	}

	var promotion PieceType = None
	if san[len(san)-2] == '=' {
		p, _, valid := pieceFromFenChar(san[len(san)-1])
		if !valid {
			return Move(0), fmt.Errorf("%w: %s", ErrInvalidPromotion, originalSan)
		}
		promotion = p
		san = san[:len(san)-2]
	} else if len(san) >= 3 && isPieceChar(san[len(san)-1]) && !isRank(san[len(san)-1]) {
		p, _, valid := pieceFromFenChar(san[len(san)-1])
		if valid {
			promotion = p
			san = san[:len(san)-1]
		}
	}

	if len(san) < 2 {
		return Move(0), fmt.Errorf("%w: %s", ErrSANTooShort, originalSan)
	}

	destFile := san[len(san)-2]
	destRank := san[len(san)-1]
	if !isFile(destFile) || !isRank(destRank) {
		return Move(0), fmt.Errorf("%w: %s", ErrInvalidDestination, originalSan)
	}
	toSquare := getNumBySquareName(destFile, destRank)
	san = san[:len(san)-2]

	var flags uint16 = FlagQuiet

	if len(san) > 0 && san[len(san)-1] == 'x' {
		flags = FlagCapture
		san = san[:len(san)-1]
	}

	if promotion != None {
		isCapture := (flags == FlagCapture)
		switch promotion {
		case Knight:
			flags = FlagPromoKnight
		case Bishop:
			flags = FlagPromoBishop
		case Rook:
			flags = FlagPromoRook
		case Queen:
			flags = FlagPromoQueen
		default:
			return Move(0), fmt.Errorf("%w: %s", ErrInvalidPromotion, originalSan)
		}
		if isCapture {
			flags |= FlagCapture
		}
	}

	piece := Pawn
	var fileHint, rankHint byte

	if len(san) > 0 {
		firstChar := san[0]
		if isPieceChar(firstChar) {
			p, _, _ := pieceFromFenChar(firstChar)
			piece = p
			san = san[1:]
		}
	}

	// Any remaining characters represent disambiguation (e.g., 'b'  in Nbd7, '1' in R1xd1)
	if len(san) > 0 {
		if len(san) == 1 {
			if isFile(san[0]) {
				fileHint = san[0]
			} else if isRank(san[0]) {
				rankHint = san[0]
			} else {
				return Move(0), fmt.Errorf("%w: %s", ErrInvalidDisambiguation, originalSan)
			}
		} else if len(san) == 2 {
			if isFile(san[0]) && isRank(san[1]) {
				fileHint = san[0]
				rankHint = san[1]
			} else {
				return Move(0), fmt.Errorf("%w: %s", ErrInvalidDoubleDisambiguation, originalSan)
			}
		} else {
			return Move(0), fmt.Errorf("%w: %s", ErrInvalidDisambiguation, originalSan)
		}
	}

	fromSquare := findFromSquare(pos, piece, color, toSquare, fileHint, rankHint)
	if fromSquare == NoSquare {
		return Move(0), fmt.Errorf("%w: %s", ErrIllegalOrAmbiguousMove, originalSan)
	}

	return NewMove(fromSquare, toSquare, flags), nil
}

// isFile checks if a byte represents a valid file (a-h).
func isFile(b byte) bool {
	return b >= 'a' && b <= 'h'
}

// isRank checks if a byte represents a valid rank (1-8).
func isRank(b byte) bool {
	return b >= '1' && b <= '8'
}

// isPieceChar checks if a byte is a valid piece identifier.
func isPieceChar(b byte) bool {
	return b == 'K' || b == 'Q' || b == 'R' || b == 'B' || b == 'N'
}

// getNumBySquareName converts file and rank bytes into a 0-63 square index.
func getNumBySquareName(file, rank byte) Square {
	f := int(file - 'a')
	r := int('8' - rank)
	return Square(r*8 + f)
}

// findFromSquare finds the origin square for a piece moving to the given destination.
func findFromSquare(pos *Position, piece PieceType, color PieceColor, toSquare Square, fileHint, rankHint byte) Square {
	bb := uint64(pos.Pieces[int(piece)+int(color)])
	if bb == 0 {
		return NoSquare
	}

	var own, enemy Bitboard
	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	for pt := Pawn; pt <= King; pt++ {
		own |= pos.Pieces[int(pt)+int(color)]
		enemy |= pos.Pieces[int(pt)+int(enemyColor)]
	}

	toMask := Bitboard(1 << uint64(toSquare))
	epMask := Bitboard(0)
	if pos.EnPassantSquare != NoSquare {
		epMask = Bitboard(1 << uint64(pos.EnPassantSquare))
	}

	for bb != 0 {
		sq := Square(bits.TrailingZeros64(bb))
		bb &= bb - 1

		if fileHint != 0 && (sq%8) != Square(fileHint-'a') {
			continue
		}
		if rankHint != 0 && (sq/8) != Square('8'-rankHint) {
			continue
		}

		pieceMask := Bitboard(1 << uint64(sq))
		var moves Bitboard

		switch piece {
		case Pawn:
			moves = GeneratePawnMoves(pieceMask, own, enemy, color, epMask)
		case Knight:
			moves = GenerateKnightMoves(pieceMask, own)
		case Bishop:
			moves = GenerateBishopMoves(pieceMask, own, enemy)
		case Rook:
			moves = GenerateRookMoves(pieceMask, own, enemy)
		case Queen:
			moves = GenerateQueenMoves(pieceMask, own, enemy)
		case King:
			moves = GenerateKingMoves(pieceMask, own, enemy)
		}

		if (moves & toMask) != 0 {
			return sq
		}
	}

	return NoSquare
}
