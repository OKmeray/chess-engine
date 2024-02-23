from logic.bitboard import Position, PieceColor


class BoardState:
    def __init__(self, position: Position, side_to_move: PieceColor = PieceColor.WHITE):
        self.position = position
        self.side_to_move = side_to_move
        self.evaluation = 0

    # 4 functions
