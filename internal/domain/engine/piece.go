package engine

// PieceColor represents the color of a chess piece or side to move.
type PieceColor int

// Side to move and piece colors.
const (
	Black PieceColor = 0
	White PieceColor = 6
)

// PieceType represents a piece type on the board.
type PieceType int

// Piece type identifiers.
const (
	None   PieceType = -1
	Pawn   PieceType = 0
	Knight PieceType = 1
	Bishop PieceType = 2
	Rook   PieceType = 3
	Queen  PieceType = 4
	King   PieceType = 5
)

// PiecePrice represents the relative material value of a piece in centipawns.
type PiecePrice int

// Piece values in centipawns.
const (
	PawnPrice   PiecePrice = 100
	KnightPrice PiecePrice = 300
	BishopPrice PiecePrice = 300
	RookPrice   PiecePrice = 500
	QueenPrice  PiecePrice = 900
	KingPrice   PiecePrice = 20000
)
