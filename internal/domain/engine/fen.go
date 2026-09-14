package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyFEN            = errors.New("FEN string cannot be empty")
	ErrInvalidFENFormat    = errors.New("invalid FEN string format")
	ErrInvalidFENPiece     = errors.New("invalid FEN piece character")
	ErrInvalidFENLength    = errors.New("invalid FEN piece placement length")
	ErrInvalidFENSide      = errors.New("invalid FEN side to move")
	ErrInvalidFENCastling  = errors.New("invalid FEN castling right")
	ErrInvalidFENEnPassant = errors.New("invalid FEN en passant square")
	ErrMissingKing         = errors.New("invalid position: missing king")
)

// ParseFEN initializes a new Position from a FEN string.
func ParseFEN(fen string) (*Position, error) {
	if strings.TrimSpace(fen) == "" {
		return nil, ErrEmptyFEN
	}

	parts := strings.Fields(fen)
	if len(parts) != 6 {
		return nil, ErrInvalidFENFormat
	}

	pos := &Position{}

	if err := parsePiecePlacement(pos, parts[0]); err != nil {
		return nil, err
	}
	if err := parseSideToMove(pos, parts[1]); err != nil {
		return nil, err
	}
	if err := parseCastlingRights(pos, parts[2]); err != nil {
		return nil, err
	}
	if err := parseEnPassant(pos, parts[3]); err != nil {
		return nil, err
	}

	halfMoves, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, fmt.Errorf("invalid halfmove clock: %w", err)
	}
	pos.HalfMoves = halfMoves

	currentTurn, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, fmt.Errorf("invalid fullmove number: %w", err)
	}
	pos.CurrentTurn = currentTurn

	// Ensure ByColor is populated
	for pt := Pawn; pt <= King; pt++ {
		pos.ByColor[1] |= pos.Pieces[int(pt)+int(White)]
		pos.ByColor[0] |= pos.Pieces[int(pt)+int(Black)]
	}

	// Validate that both kings exist
	if pos.Pieces[int(King)+int(White)] == 0 || pos.Pieces[int(King)+int(Black)] == 0 {
		return nil, ErrMissingKing
	}

	return pos, nil
}

// FEN generates the FEN string representation of the position.
func (p *Position) FEN() string {
	if p == nil {
		return ""
	}

	var sb strings.Builder
	sb.Grow(80)
	emptyCount := 0

	for i := 0; i < 64; i++ {
		piece, color := p.GetPieceAndColorBySquare(Square(i))

		if piece == None {
			emptyCount++
		} else {
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
				emptyCount = 0
			}
			sb.WriteByte(pieceToFenChar(piece, color))
		}

		if (i+1)%8 == 0 {
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
				emptyCount = 0
			}
			if i != 63 {
				sb.WriteByte('/')
			}
		}
	}

	sb.WriteByte(' ')

	if p.SideToMove == White {
		sb.WriteByte('w')
	} else {
		sb.WriteByte('b')
	}

	sb.WriteByte(' ')

	anyCastle := false
	if p.CastlingRights&WhiteShort == WhiteShort {
		sb.WriteByte('K')
		anyCastle = true
	}
	if p.CastlingRights&WhiteLong == WhiteLong {
		sb.WriteByte('Q')
		anyCastle = true
	}
	if p.CastlingRights&BlackShort == BlackShort {
		sb.WriteByte('k')
		anyCastle = true
	}
	if p.CastlingRights&BlackLong == BlackLong {
		sb.WriteByte('q')
		anyCastle = true
	}
	if !anyCastle {
		sb.WriteByte('-')
	}

	sb.WriteByte(' ')

	if p.EnPassantSquare != NoSquare {
		sb.WriteString(GetSquareNameByNum(p.EnPassantSquare))
	} else {
		sb.WriteByte('-')
	}

	sb.WriteByte(' ')

	sb.WriteString(strconv.Itoa(p.HalfMoves))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(p.CurrentTurn))

	return sb.String()
}

// GetSquareNameByNum converts a Square index to its algebraic notation.
func GetSquareNameByNum(num Square) string {
	if num == NoSquare {
		return "-"
	}

	file := byte(num%8) + 'a'
	rank := '8' - byte(num/8)

	return string([]byte{file, rank})
}

// GetNumBySquareName converts an algebraic notation string to a Square.
func GetNumBySquareName(square string) Square {
	if len(square) != 2 {
		return NoSquare
	}

	if square[0] < 'a' || square[0] > 'h' || square[1] < '1' || square[1] > '8' {
		return NoSquare
	}

	f := int(square[0] - 'a')
	r := int('8' - square[1])

	return Square(r*8 + f)
}

// parsePiecePlacement parses the piece placement field of a FEN string.
func parsePiecePlacement(pos *Position, placement string) error {
	square := Square(0)
	for i := 0; i < len(placement); i++ {
		c := placement[i]
		if c >= '1' && c <= '8' {
			square += Square(c - '0')
		} else if c != '/' {
			piece, color, valid := pieceFromFenChar(c)
			if !valid {
				return fmt.Errorf("%w: %c", ErrInvalidFENPiece, c)
			}

			pos.AddPiece(piece, color, square)
			square++
		}
	}

	if square != 64 {
		return fmt.Errorf("%w: expected 64 squares, got %d", ErrInvalidFENLength, square)
	}
	return nil
}

// parseSideToMove parses the active color field of a FEN string.
func parseSideToMove(pos *Position, side string) error {
	switch side {
	case "w":
		pos.SideToMove = White
	case "b":
		pos.SideToMove = Black
	default:
		return fmt.Errorf("%w: %s", ErrInvalidFENSide, side)
	}
	return nil
}

// parseCastlingRights parses the castling availability field of a FEN string.
func parseCastlingRights(pos *Position, rights string) error {
	pos.CastlingRights = 0

	if rights == "-" {
		return nil
	}

	for i := 0; i < len(rights); i++ {
		switch rights[i] {
		case 'K':
			pos.CastlingRights |= WhiteShort
		case 'Q':
			pos.CastlingRights |= WhiteLong
		case 'k':
			pos.CastlingRights |= BlackShort
		case 'q':
			pos.CastlingRights |= BlackLong
		default:
			return fmt.Errorf("%w: %c", ErrInvalidFENCastling, rights[i])
		}
	}
	return nil
}

// parseEnPassant parses the en passant target square field of a FEN string.
func parseEnPassant(pos *Position, ep string) error {
	if ep != "-" {
		sq := GetNumBySquareName(ep)
		if sq == NoSquare {
			return ErrInvalidFENEnPassant
		}
		pos.EnPassantSquare = sq
	} else {
		pos.EnPassantSquare = NoSquare
	}
	return nil
}

// pieceFromFenChar maps a FEN character to its PieceType and PieceColor.
func pieceFromFenChar(c byte) (PieceType, PieceColor, bool) {
	switch c {
	case 'p':
		return Pawn, Black, true
	case 'n':
		return Knight, Black, true
	case 'b':
		return Bishop, Black, true
	case 'r':
		return Rook, Black, true
	case 'q':
		return Queen, Black, true
	case 'k':
		return King, Black, true
	case 'P':
		return Pawn, White, true
	case 'N':
		return Knight, White, true
	case 'B':
		return Bishop, White, true
	case 'R':
		return Rook, White, true
	case 'Q':
		return Queen, White, true
	case 'K':
		return King, White, true
	default:
		return 0, 0, false
	}
}

// pieceToFenChar converts a PieceType and PieceColor into a FEN byte.
func pieceToFenChar(p PieceType, c PieceColor) byte {
	var char byte

	switch p {
	case Pawn:
		char = 'p'
	case Knight:
		char = 'n'
	case Bishop:
		char = 'b'
	case Rook:
		char = 'r'
	case Queen:
		char = 'q'
	case King:
		char = 'k'
	default:
		return '?'
	}

	if c == White {
		char -= 32 // shift in ASCII table from lowercase letter to uppercase
	}

	return char
}
