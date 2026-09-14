package engine

import "math/bits"

// generatePseudoLegalMoves generates all pseudo-legal moves for the current side to move.
func (p *Position) generatePseudoLegalMoves(moves []Move) []Move {
	color := p.SideToMove
	enemyColor := 6 - color
	own := p.ByColor[color/6]
	enemy := p.ByColor[enemyColor/6]

	epMask := Bitboard(0)
	if p.EnPassantSquare != NoSquare {
		epMask = Bitboard(1 << uint64(p.EnPassantSquare))
	}

	for pt := Pawn; pt <= King; pt++ {
		bb := p.Pieces[int(pt)+int(color)]
		for bb != 0 {
			fromSq := Square(bits.TrailingZeros64(uint64(bb)))
			bb &= bb - 1

			pieceMask := Bitboard(1 << uint64(fromSq))
			var targets Bitboard

			switch pt {
			case Pawn:
				targets = GeneratePawnMoves(pieceMask, own, enemy, color, epMask)
			case Knight:
				targets = GenerateKnightMoves(pieceMask, own)
			case Bishop:
				targets = GenerateBishopMoves(pieceMask, own, enemy)
			case Rook:
				targets = GenerateRookMoves(pieceMask, own, enemy)
			case Queen:
				targets = GenerateQueenMoves(pieceMask, own, enemy)
			case King:
				targets = GenerateKingMoves(pieceMask, own, enemy)
			}

			for targets != 0 {
				toSq := Square(bits.TrailingZeros64(uint64(targets)))
				targets &= targets - 1

				flags := FlagQuiet
				if (Bitboard(1<<uint64(toSq)) & enemy) != 0 {
					flags = FlagCapture
				}

				if pt == Pawn {
					// Double pawn push
					if (color == White && fromSq/8 == 6 && toSq/8 == 4) ||
						(color == Black && fromSq/8 == 1 && toSq/8 == 3) {
						flags = FlagDoublePawn
					} else if toSq == p.EnPassantSquare {
						flags = FlagEnPassant
					} else if (color == White && toSq/8 == 0) || (color == Black && toSq/8 == 7) {
						// Promotion
						promos := []uint16{FlagPromoQueen, FlagPromoKnight, FlagPromoRook, FlagPromoBishop}
						for _, promo := range promos {
							f := promo
							if flags == FlagCapture {
								f |= FlagCapture
							}
							moves = append(moves, NewMove(fromSq, toSq, f))
						}
						continue
					}
				}

				moves = append(moves, NewMove(fromSq, toSq, flags))
			}
		}
	}

	// Castling
	if color == White {
		if p.CastlingRights.Has(WhiteShort) && (own|enemy)&(1<<61|1<<62) == 0 {
			if !p.IsSquareAttacked(60, Black) && !p.IsSquareAttacked(61, Black) {
				moves = append(moves, NewMove(60, 62, FlagKingCastle))
			}
		}
		if p.CastlingRights.Has(WhiteLong) && (own|enemy)&(1<<57|1<<58|1<<59) == 0 {
			if !p.IsSquareAttacked(60, Black) && !p.IsSquareAttacked(59, Black) {
				moves = append(moves, NewMove(60, 58, FlagQueenCastle))
			}
		}
	} else {
		if p.CastlingRights.Has(BlackShort) && (own|enemy)&(1<<5|1<<6) == 0 {
			if !p.IsSquareAttacked(4, White) && !p.IsSquareAttacked(5, White) {
				moves = append(moves, NewMove(4, 6, FlagKingCastle))
			}
		}
		if p.CastlingRights.Has(BlackLong) && (own|enemy)&(1<<1|1<<2|1<<3) == 0 {
			if !p.IsSquareAttacked(4, White) && !p.IsSquareAttacked(3, White) {
				moves = append(moves, NewMove(4, 2, FlagQueenCastle))
			}
		}
	}

	return moves
}

// GenerateMoves generates all legal moves for the current side to move.
func (p *Position) GenerateMoves(moves []Move) []Move {
	var pseudoBuf [256]Move
	pseudoMoves := p.generatePseudoLegalMoves(pseudoBuf[:0])

	for _, move := range pseudoMoves {
		undo := p.MakeMove(move)
		sideThatMoved := 6 - p.SideToMove
		kingSq := p.GetKingSq(sideThatMoved)

		// If the king is not attacked by the opponent after the move, it's legal
		if !p.IsSquareAttacked(kingSq, p.SideToMove) {
			moves = append(moves, move)
		}
		p.UnmakeMove(move, undo)
	}
	return moves
}
