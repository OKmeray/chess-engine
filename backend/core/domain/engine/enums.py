from enum import IntEnum


# # ------------------------
# WHITE_PAWN = 0b1001  # 9
# WHITE_KNIGHT = 0b1010  # 10
# WHITE_BISHOP = 0b1011  # 11
# WHITE_ROOK = 0b1100  # 12
# WHITE_QUEEN = 0b1101  # 13
# WHITE_KING = 0b1110  # 14
# # ------------------------
# BLACK_PAWN = 0b0001  # 1
# BLACK_KNIGHT = 0b0010  # 2
# BLACK_BISHOP = 0b0011  # 3
# BLACK_ROOK = 0b0100  # 4
# BLACK_QUEEN = 0b0101  # 5
# BLACK_KING = 0b0110  # 6
# # ------------------------
# <<  ->  +
# >>  ->  -


class PieceEnum(IntEnum):
    NONE = 0
    PAWN = 1
    KNIGHT = 2
    BISHOP = 3
    ROOK = 4
    QUEEN = 5
    KING = 6


class PieceColor(IntEnum):
    BLACK = 0
    WHITE = 8


class CastleEnum:
    BlackShortCastle = 0
    BlackLongCastle = 1
    WhiteShortCastle = 8
    WhiteLongCastle = 9


class PiecePrice(IntEnum):
    PAWN = 100
    KNIGHT = 300  # 320
    BISHOP = 300  # 330
    ROOK = 500
    QUEEN = 900
    KING = 20000


class PositionOutcome(IntEnum):
    LOSS = 0
    DRAW = 1
    WIN = 2
