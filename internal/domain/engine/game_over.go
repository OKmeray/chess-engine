package engine

// IsCheckmate returns true if the current side to move has no legal moves and is in check.
func (p *Position) IsCheckmate() bool {
	legalMoves := p.GenerateMoves(nil)
	if len(legalMoves) > 0 {
		return false
	}
	return p.IsSquareAttacked(p.GetKingSq(p.SideToMove), 6-p.SideToMove)
}

// IsStalemate returns true if the current side to move has no legal moves but is not in check.
func (p *Position) IsStalemate() bool {
	legalMoves := p.GenerateMoves(nil)
	if len(legalMoves) > 0 {
		return false
	}
	return !p.IsSquareAttacked(p.GetKingSq(p.SideToMove), 6-p.SideToMove)
}

// IsFiftyMoveRule returns true if 100 half-moves have been played without a capture or pawn push.
func (p *Position) IsFiftyMoveRule() bool {
	return p.HalfMoves >= 100
}

// IsInsufficientMaterial returns true if neither side has enough pieces to force a checkmate.
func (p *Position) IsInsufficientMaterial() bool {
	// First, check if there are any pawns, rooks, or queens. If so, material is sufficient.
	if p.Pieces[int(White)+int(Pawn)] != 0 || p.Pieces[int(Black)+int(Pawn)] != 0 ||
		p.Pieces[int(White)+int(Rook)] != 0 || p.Pieces[int(Black)+int(Rook)] != 0 ||
		p.Pieces[int(White)+int(Queen)] != 0 || p.Pieces[int(Black)+int(Queen)] != 0 {
		return false
	}

	whiteKnights := p.Pieces[int(White)+int(Knight)].Count()
	blackKnights := p.Pieces[int(Black)+int(Knight)].Count()
	totalKnights := whiteKnights + blackKnights

	whiteBishops := p.Pieces[int(White)+int(Bishop)]
	blackBishops := p.Pieces[int(Black)+int(Bishop)]
	totalBishopsCount := whiteBishops.Count() + blackBishops.Count()

	totalMinorPieces := totalKnights + totalBishopsCount

	// King vs King, King and Knight vs King, or King and Bishop vs King
	if totalMinorPieces <= 1 {
		return true
	}

	// King and Bishop(s) vs King and Bishop(s) of the SAME color complex
	// (Any number of bishops on the same color is insufficient to force mate)
	if totalKnights == 0 {
		allBishops := whiteBishops | blackBishops

		const darkSquares Bitboard = 0b10101010_01010101_10101010_01010101_10101010_01010101_10101010_01010101
		const lightSquares Bitboard = 0b01010101_10101010_01010101_10101010_01010101_10101010_01010101_10101010

		if allBishops&darkSquares == allBishops || allBishops&lightSquares == allBishops {
			return true
		}
	}

	return false
}

// IsThreeFoldRepetition returns true if the current position has occurred 2+ times before.
func (p *Position) IsThreeFoldRepetition() bool {
	// We only need to search up to p.HalfMoves, because any capture or pawn push
	// irrevocably changes the board and resets HalfMoves to 0.
	count := 0
	for i := range p.HalfMoves {
		if p.RepetitionHistory[i] == p.Hash {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}
