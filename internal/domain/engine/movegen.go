package engine

// Constants representing ranks and files
const (
	Rank1 Bitboard = 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	Rank2 Bitboard = 0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000
	Rank7 Bitboard = 0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000
	Rank8 Bitboard = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111

	FileA Bitboard = 0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001
	FileB Bitboard = 0b00000010_00000010_00000010_00000010_00000010_00000010_00000010_00000010
	FileE Bitboard = 0b00010000_00010000_00010000_00010000_00010000_00010000_00010000_00010000
	FileG Bitboard = 0b01000000_01000000_01000000_01000000_01000000_01000000_01000000_01000000
	FileH Bitboard = 0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000
)

// GeneratePawnMoves generates a bitboard mask of all valid pawn moves
func GeneratePawnMoves(piece, own, enemy Bitboard, color PieceColor, epSquare Bitboard) Bitboard {
	var moves Bitboard

	if color == White {
		// en passant and taking diagonal
		target := piece >> 9
		if (piece&FileA) == 0 && ((target&enemy) != 0 || target == epSquare) {
			moves |= target
		}

		// en passant and taking diagonal
		target = piece >> 7
		if (piece&FileH) == 0 && ((target&enemy) != 0 || target == epSquare) {
			moves |= target
		}

		// move forward 1 square or 2 squares
		target = piece >> 8
		if (target & (enemy | own)) == 0 {
			moves |= target
			if (piece & Rank2) != 0 {
				doubleForward := piece >> 16
				if (doubleForward & (enemy | own)) == 0 {
					moves |= doubleForward
				}
			}
		}
	} else {
		// en passant and taking diagonal
		target := piece << 9
		if (piece&FileH) == 0 && ((target&enemy) != 0 || target == epSquare) {
			moves |= target
		}

		// en passant and taking diagonal
		target = piece << 7
		if (piece&FileA) == 0 && ((target&enemy) != 0 || target == epSquare) {
			moves |= target
		}

		// move forward 1 square or 2 squares
		target = piece << 8
		if (target & (enemy | own)) == 0 {
			moves |= target
			if (piece & Rank7) != 0 {
				doubleForward := piece << 16
				if (doubleForward & (enemy | own)) == 0 {
					moves |= doubleForward
				}
			}
		}
	}

	return moves
}

// GenerateKnightMoves generates a bitboard mask of all valid knight moves
func GenerateKnightMoves(piece, own Bitboard) Bitboard {
	var moves Bitboard

	// Upward moves
	if target := (piece >> 17) &^ FileH; target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece >> 15) &^ FileA; target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece >> 10) &^ (FileG | FileH); target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece >> 6) &^ (FileA | FileB); target != 0 && (target&own) == 0 {
		moves |= target
	}

	// Downward moves
	if target := (piece << 17) &^ FileA; target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece << 15) &^ FileH; target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece << 10) &^ (FileA | FileB); target != 0 && (target&own) == 0 {
		moves |= target
	}
	if target := (piece << 6) &^ (FileG | FileH); target != 0 && (target&own) == 0 {
		moves |= target
	}

	return moves
}

// GenerateBishopMoves generates a bitboard mask of all valid bishop moves
func GenerateBishopMoves(piece, own, enemy Bitboard) Bitboard {
	return generateDiagonalUp(piece, own, enemy) |
		generateDiagonalDown(piece, own, enemy) |
		generateAntidiagonalUp(piece, own, enemy) |
		generateAntidiagonalDown(piece, own, enemy)
}

// GenerateRookMoves generates a bitboard mask of all valid rook moves
func GenerateRookMoves(piece, own, enemy Bitboard) Bitboard {
	return generateVerticalUp(piece, own, enemy) |
		generateVerticalDown(piece, own, enemy) |
		generateHorizontalLeft(piece, own, enemy) |
		generateHorizontalRight(piece, own, enemy)
}

// GenerateQueenMoves generates a bitboard mask of all valid queen moves
func GenerateQueenMoves(piece, own, enemy Bitboard) Bitboard {
	return GenerateRookMoves(piece, own, enemy) |
		GenerateBishopMoves(piece, own, enemy)
}

// GenerateKingMoves generates a bitboard mask of all valid (except checks) king moves
func GenerateKingMoves(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard

	moves |= (piece >> 8)          // Up
	moves |= (piece << 8)          // Down
	moves |= (piece >> 1) &^ FileH // Left
	moves |= (piece << 1) &^ FileA // Right

	moves |= (piece >> 9) &^ FileH // Up-Left
	moves |= (piece >> 7) &^ FileA // Up-Right
	moves |= (piece << 7) &^ FileH // Down-Left
	moves |= (piece << 9) &^ FileA // Down-Right

	// Filter out squares occupied by own pieces
	return moves &^ own
}

func generateDiagonalUp(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & (FileA | Rank8)) != 0 {
			break
		}

		piece >>= 9
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateDiagonalDown(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & (FileH | Rank1)) != 0 {
			break
		}

		piece <<= 9
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateAntidiagonalUp(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & (FileH | Rank8)) != 0 {
			break
		}

		piece >>= 7
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateAntidiagonalDown(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & (FileA | Rank1)) != 0 {
			break
		}

		piece <<= 7
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateVerticalUp(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & Rank8) != 0 {
			break
		}

		piece >>= 8
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateVerticalDown(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & Rank1) != 0 {
			break
		}

		piece <<= 8
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateHorizontalLeft(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & FileA) != 0 {
			break
		}

		piece >>= 1
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}

func generateHorizontalRight(piece, own, enemy Bitboard) Bitboard {
	var moves Bitboard
	for {
		if (piece & FileH) != 0 {
			break
		}

		piece <<= 1
		if (piece & own) != 0 {
			break
		}

		moves |= piece
		if (piece & enemy) != 0 {
			break
		}
	}
	return moves
}
