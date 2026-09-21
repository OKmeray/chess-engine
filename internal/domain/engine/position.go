package engine

import "math/bits"

const bitboardCount = 12

// Position represents the state of a chess board and game.
type Position struct {
	Pieces            [bitboardCount]Bitboard
	ByColor           [2]Bitboard
	SideToMove        PieceColor
	CastlingRights    CastlingRights
	EnPassantSquare   Square
	HalfMoves         int
	CurrentTurn       int
	Hash              uint64
	RepetitionHistory [100]uint64
}

// NewPosition creates and returns a new chess position with default starting values.
func NewPosition() *Position {
	var pieces [bitboardCount]Bitboard
	for i := range bitboardCount {
		pieces[i] = NewBitboard(0)
	}
	var byColor [2]Bitboard
	for i := range 2 {
		byColor[i] = NewBitboard(0)
	}

	var sideToMove PieceColor = White

	var castlingRights CastlingRights
	castlingRights |= ^WhiteShort
	castlingRights |= ^WhiteLong
	castlingRights |= ^BlackShort
	castlingRights |= ^BlackLong

	enPassantSquare := NoSquare
	halfMoves := 0
	currentTurn := 1

	return &Position{
		Pieces:          pieces,
		ByColor:         byColor,
		SideToMove:      sideToMove,
		CastlingRights:  castlingRights,
		EnPassantSquare: enPassantSquare,
		HalfMoves:       halfMoves,
		CurrentTurn:     currentTurn,
	}
}

// AddPiece places a piece of the specified type and color on the given square.
func (p *Position) AddPiece(piece PieceType, color PieceColor, square Square) {
	indx := int(piece) + int(color)
	p.Pieces[indx].AddPiece(square)
}

// GetPieceAndColorBySquare returns the piece type and color located at the specified square.
func (p *Position) GetPieceAndColorBySquare(square Square) (PieceType, PieceColor) {
	mask := Bitboard(1 << uint64(square))
	colors := []PieceColor{White, Black}
	for _, color := range colors {
		for pt := Pawn; pt <= King; pt++ {
			if (p.Pieces[int(pt)+int(color)] & mask) != 0 {
				return pt, color
			}
		}
	}
	return None, White
}

// UndoInfo contains the information necessary to revert a move and restore the previous position state.
type UndoInfo struct {
	CapturedPiece   PieceType
	CastlingRights  CastlingRights
	EnPassantSquare Square
	HalfMoves       int
	Hash            uint64
	RepetitionHash  uint64
}

// MakeMove applies a move to the position and returns UndoInfo required to unmake it.
func (p *Position) MakeMove(move Move) UndoInfo {
	undo := UndoInfo{
		CapturedPiece:   None,
		CastlingRights:  p.CastlingRights,
		EnPassantSquare: p.EnPassantSquare,
		HalfMoves:       p.HalfMoves,
		Hash:            p.Hash,
		RepetitionHash:  p.RepetitionHistory[p.HalfMoves],
	}

	p.RepetitionHistory[p.HalfMoves] = p.Hash

	from := move.From()
	to := move.To()
	flags := move.Flags()

	piece, color := p.GetPieceAndColorBySquare(from)
	cIdx := color / 6

	// XOR out old castling and en passant
	p.Hash ^= ZobristCastling[p.CastlingRights]
	if p.EnPassantSquare != NoSquare {
		p.Hash ^= ZobristEnPassant[p.EnPassantSquare%8]
	}

	// Handle capture
	if flags == FlagCapture || flags >= FlagPromoKnightCapture {
		capturedPiece, capturedColor := p.GetPieceAndColorBySquare(to)
		undo.CapturedPiece = capturedPiece
		p.Pieces[int(capturedPiece)+int(capturedColor)] &^= Bitboard(1 << to)
		p.ByColor[capturedColor/6] &^= Bitboard(1 << to) // /6 converts White(6)->1, Black(0)->0
		p.Hash ^= ZobristTable[capturedColor/6][capturedPiece][to]

		// If a rook is captured on its starting square, revoke its owner's castling right
		if capturedPiece == Rook {
			if to == 0 {
				p.CastlingRights &^= BlackLong
			}
			if to == 7 {
				p.CastlingRights &^= BlackShort
			}
			if to == 56 {
				p.CastlingRights &^= WhiteLong
			}
			if to == 63 {
				p.CastlingRights &^= WhiteShort
			}
		}
	} else if flags == FlagEnPassant {
		undo.CapturedPiece = Pawn
		capSq := to + 8
		if color == Black {
			capSq = to - 8
		}
		p.Pieces[int(Pawn)+int(6-color)] &^= Bitboard(1 << capSq)
		p.ByColor[(6-color)/6] &^= Bitboard(1 << capSq)
		p.Hash ^= ZobristTable[(6-color)/6][Pawn][capSq]
	}

	// Move the piece
	p.Pieces[int(piece)+int(color)] &^= Bitboard(1 << from)
	p.Pieces[int(piece)+int(color)] |= Bitboard(1 << to)
	p.ByColor[cIdx] &^= Bitboard(1 << from)
	p.ByColor[cIdx] |= Bitboard(1 << to)

	p.Hash ^= ZobristTable[cIdx][piece][from]
	p.Hash ^= ZobristTable[cIdx][piece][to]

	// Handle Castling moves
	switch flags {
	case FlagKingCastle:
		rookFrom, rookTo := Square(from+3), Square(to-1)
		p.Pieces[int(Rook)+int(color)] &^= Bitboard(1 << rookFrom)
		p.Pieces[int(Rook)+int(color)] |= Bitboard(1 << rookTo)
		p.ByColor[cIdx] &^= Bitboard(1 << rookFrom)
		p.ByColor[cIdx] |= Bitboard(1 << rookTo)

		p.Hash ^= ZobristTable[cIdx][Rook][rookFrom]
		p.Hash ^= ZobristTable[cIdx][Rook][rookTo]

	case FlagQueenCastle:
		rookFrom, rookTo := Square(from-4), Square(to+1)
		p.Pieces[int(Rook)+int(color)] &^= Bitboard(1 << rookFrom)
		p.Pieces[int(Rook)+int(color)] |= Bitboard(1 << rookTo)
		p.ByColor[cIdx] &^= Bitboard(1 << rookFrom)
		p.ByColor[cIdx] |= Bitboard(1 << rookTo)

		p.Hash ^= ZobristTable[cIdx][Rook][rookFrom]
		p.Hash ^= ZobristTable[cIdx][Rook][rookTo]
	}

	// Handle Promotion
	if promoType := move.Promotion(); promoType != None {
		p.Pieces[int(Pawn)+int(color)] &^= Bitboard(1 << to)
		p.Pieces[int(promoType)+int(color)] |= Bitboard(1 << to)

		p.Hash ^= ZobristTable[cIdx][Pawn][to]
		p.Hash ^= ZobristTable[cIdx][promoType][to]
	}

	// Update Castling Rights
	switch piece {
	case King:
		if color == White {
			p.CastlingRights &^= (WhiteShort | WhiteLong)
		} else {
			p.CastlingRights &^= (BlackShort | BlackLong)
		}
	case Rook:
		if from == 0 {
			p.CastlingRights &^= BlackLong
		}
		if from == 7 {
			p.CastlingRights &^= BlackShort
		}
		if from == 56 {
			p.CastlingRights &^= WhiteLong
		}
		if from == 63 {
			p.CastlingRights &^= WhiteShort
		}
	}

	// Update En Passant Square
	p.EnPassantSquare = NoSquare
	if flags == FlagDoublePawn {
		if color == White {
			p.EnPassantSquare = to + 8
		} else {
			p.EnPassantSquare = to - 8
		}
	}

	// Update HalfMoves
	if piece == Pawn || undo.CapturedPiece != None {
		p.HalfMoves = 0
	} else {
		p.HalfMoves++
	}

	if p.SideToMove == Black {
		p.CurrentTurn++
	}
	p.SideToMove = 6 - p.SideToMove // Switch color

	// XOR in new castling, en passant, and side to move
	p.Hash ^= ZobristCastling[p.CastlingRights]
	if p.EnPassantSquare != NoSquare {
		p.Hash ^= ZobristEnPassant[p.EnPassantSquare%8]
	}
	p.Hash ^= ZobristSideToMove

	return undo
}

// UnmakeMove reverts the state of the position using UndoInfo.
func (p *Position) UnmakeMove(move Move, undo UndoInfo) {
	p.RepetitionHistory[undo.HalfMoves] = undo.RepetitionHash
	p.Hash = undo.Hash
	p.SideToMove = 6 - p.SideToMove
	if p.SideToMove == Black {
		p.CurrentTurn--
	}

	from := move.From()
	to := move.To()
	flags := move.Flags()
	color := p.SideToMove

	piece := Pawn
	if promoType := move.Promotion(); promoType != None {
		// Demote
		p.Pieces[int(promoType)+int(color)] &^= Bitboard(1 << to)
		p.Pieces[int(Pawn)+int(color)] |= Bitboard(1 << to)
	} else {
		piece, _ = p.GetPieceAndColorBySquare(to)
	}

	// Move the piece back
	p.Pieces[int(piece)+int(color)] &^= Bitboard(1 << to)
	p.Pieces[int(piece)+int(color)] |= Bitboard(1 << from)
	p.ByColor[color/6] &^= Bitboard(1 << to)
	p.ByColor[color/6] |= Bitboard(1 << from)

	// Restore Castling
	switch flags {
	case FlagKingCastle:
		rookFrom, rookTo := Square(from+3), Square(to-1)
		p.Pieces[int(Rook)+int(color)] &^= Bitboard(1 << rookTo)
		p.Pieces[int(Rook)+int(color)] |= Bitboard(1 << rookFrom)
		p.ByColor[color/6] &^= Bitboard(1 << rookTo)
		p.ByColor[color/6] |= Bitboard(1 << rookFrom)
	case FlagQueenCastle:
		rookFrom, rookTo := Square(from-4), Square(to+1)
		p.Pieces[int(Rook)+int(color)] &^= Bitboard(1 << rookTo)
		p.Pieces[int(Rook)+int(color)] |= Bitboard(1 << rookFrom)
		p.ByColor[color/6] &^= Bitboard(1 << rookTo)
		p.ByColor[color/6] |= Bitboard(1 << rookFrom)
	}

	// Restore Captures
	if flags == FlagEnPassant {
		capSq := to + 8
		if color == Black {
			capSq = to - 8
		}
		p.Pieces[int(Pawn)+int(6-color)] |= Bitboard(1 << capSq)
		p.ByColor[(6-color)/6] |= Bitboard(1 << capSq)
	} else if undo.CapturedPiece != None {
		p.Pieces[int(undo.CapturedPiece)+int(6-color)] |= Bitboard(1 << to)
		p.ByColor[(6-color)/6] |= Bitboard(1 << to)
	}

	p.CastlingRights = undo.CastlingRights
	p.EnPassantSquare = undo.EnPassantSquare
	p.HalfMoves = undo.HalfMoves
}

// IsSquareAttacked determines if a specific square is attacked by pieces of a given color.
func (p *Position) IsSquareAttacked(sq Square, attackerColor PieceColor) bool {
	sqMask := Bitboard(1 << sq)
	own := p.ByColor[(6-attackerColor)/6]
	enemy := p.ByColor[attackerColor/6]

	// Check Pawns
	if attackerColor == White {
		if (sqMask&^FileH != 0) && (p.Pieces[int(Pawn)+int(White)]&(sqMask<<9)) != 0 {
			return true
		}
		if (sqMask&^FileA != 0) && (p.Pieces[int(Pawn)+int(White)]&(sqMask<<7)) != 0 {
			return true
		}
	} else {
		if (sqMask&^FileA != 0) && (p.Pieces[int(Pawn)+int(Black)]&(sqMask>>9)) != 0 {
			return true
		}
		if (sqMask&^FileH != 0) && (p.Pieces[int(Pawn)+int(Black)]&(sqMask>>7)) != 0 {
			return true
		}
	}

	// Check Knights
	if (GenerateKnightMoves(sqMask, 0) & p.Pieces[int(Knight)+int(attackerColor)]) != 0 {
		return true
	}

	// Check Kings
	if (GenerateKingMoves(sqMask, 0, 0) & p.Pieces[int(King)+int(attackerColor)]) != 0 {
		return true
	}

	// Check Bishops/Queens
	if (GenerateBishopMoves(sqMask, own, enemy) & (p.Pieces[int(Bishop)+int(attackerColor)] | p.Pieces[int(Queen)+int(attackerColor)])) != 0 {
		return true
	}

	// Check Rooks/Queens
	if (GenerateRookMoves(sqMask, own, enemy) & (p.Pieces[int(Rook)+int(attackerColor)] | p.Pieces[int(Queen)+int(attackerColor)])) != 0 {
		return true
	}

	return false
}

// GetKingSq returns the square where the given color's king is located.
func (p *Position) GetKingSq(color PieceColor) Square {
	bb := p.Pieces[int(King)+int(color)]
	return Square(bits.TrailingZeros64(uint64(bb)))
}
