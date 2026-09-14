package engine

// PieceColor
type PieceColor int

const (
	Black PieceColor = 0
	White PieceColor = 6
)

// Piece
type PieceType int

const (
	None   PieceType = -1
	Pawn   PieceType = 0
	Knight PieceType = 1
	Bishop PieceType = 2
	Rook   PieceType = 3
	Queen  PieceType = 4
	King   PieceType = 5
)

// PiecePrice
type PiecePrice int

const (
	PawnPrice   PiecePrice = 100
	KnightPrice PiecePrice = 300
	BishopPrice PiecePrice = 300
	RookPrice   PiecePrice = 500
	QueenPrice  PiecePrice = 900
	KingPrice   PiecePrice = 20000
)
