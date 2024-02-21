from enum import IntEnum

from functionality_test import get_all_moves_by_piece, position
from logic.bitboard import PieceEnum, PieceColor
from logic.helping_function import get_separate_piece, get_num_from_bitboard


class PiecePrice(IntEnum):
    PAWN = 100
    KNIGHT = 300  # 320
    BISHOP = 300  # 330
    ROOK = 500
    QUEEN = 900
    KING = 20000


pawn_square_table = [
    0,  0,  0,  0,  0,  0,  0,  0,
    50, 50, 50, 50, 50, 50, 50, 50,
    10, 10, 20, 30, 30, 20, 10, 10,
     5,  5, 10, 25, 25, 10,  5,  5,
     0,  0,  0, 20, 20,  0,  0,  0,
     5, -5,-10,  0,  0,-10, -5,  5,
     5, 10, 10,-20,-20, 10, 10,  5,
     0,  0,  0,  0,  0,  0,  0,  0
]

knight_square_table = [
    -50,-40,-30,-30,-30,-30,-40,-50,
    -40,-20,  0,  0,  0,  0,-20,-40,
    -30,  0, 10, 15, 15, 10,  0,-30,
    -30,  5, 15, 20, 20, 15,  5,-30,
    -30,  0, 15, 20, 20, 15,  0,-30,
    -30,  5, 10, 15, 15, 10,  5,-30,
    -40,-20,  0,  5,  5,  0,-20,-40,
    -50,-40,-30,-30,-30,-30,-40,-50,
]


bishop_square_table = [
    -20,-10,-10,-10,-10,-10,-10,-20,
    -10,  0,  0,  0,  0,  0,  0,-10,
    -10,  0,  5, 10, 10,  5,  0,-10,
    -10,  5,  5, 10, 10,  5,  5,-10,
    -10,  0, 10, 10, 10, 10,  0,-10,
    -10, 10, 10, 10, 10, 10, 10,-10,
    -10,  5,  0,  0,  0,  0,  5,-10,
    -20,-10,-10,-10,-10,-10,-10,-20,
]


rook_square_table = [
      0,  0,  0,  0,  0,  0,  0,  0,
      5, 10, 10, 10, 10, 10, 10,  5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
      0,  0,  0,  5,  5,  0,  0,  0
]


queen_square_table = [
    -20,-10,-10, -5, -5,-10,-10,-20,
    -10,  0,  0,  0,  0,  0,  0,-10,
    -10,  0,  5,  5,  5,  5,  0,-10,
     -5,  0,  5,  5,  5,  5,  0, -5,
      0,  0,  5,  5,  5,  5,  0, -5,
    -10,  5,  5,  5,  5,  5,  0,-10,
    -10,  0,  5,  0,  0,  0,  0,-10,
    -20,-10,-10, -5, -5,-10,-10,-20
]


king_square_middlegame_table = [
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -20,-30,-30,-40,-40,-30,-30,-20,
    -10,-20,-20,-20,-20,-20,-20,-10,
     20, 20,  0,  0,  0,  0, 20, 20,
     20, 30, 10,  0,  0, 10, 30, 20
]


king_square_endgame_table = [
    -50,-40,-30,-20,-20,-30,-40,-50,
    -30,-20,-10,  0,  0,-10,-20,-30,
    -30,-10, 20, 30, 30, 20,-10,-30,
    -30,-10, 30, 40, 40, 30,-10,-30,
    -30,-10, 30, 40, 40, 30,-10,-30,
    -30,-10, 20, 30, 30, 20,-10,-30,
    -30,-30,  0,  0,  0,  0,-30,-30,
    -50,-30,-30,-30,-30,-30,-30,-50
]


# TODO: check for mobility of pieces (by checking the amount of squares the can move, but if
## if the king is in check it should be calculated in another way)
def evaluate_position(position, side=PieceColor.WHITE):
    evaluation = 0

    evaluation += evaluate_material(position=position, side=side)
    evaluation += evaluate_piece_placement(position=position, side=side)

    return evaluation


def evaluate_material(position, side=PieceColor.WHITE):
    evaluation = 0

    for num_color in range(0, 9, 8):
        for num_piece in range(1, 7):
            separate_bit_pieces = get_separate_piece(position, PieceEnum(num_piece), PieceColor(num_color))
            separate_num_pieces = [get_num_from_bitboard(num) for num in separate_bit_pieces]
            print(PieceEnum(num_piece), PieceColor(num_color), separate_num_pieces)

            side_sign = 1 if PieceColor(num_color) == side else -1
            #for sep_num_piece in separate_num_pieces:
            for _ in separate_bit_pieces:

                if PieceEnum(num_piece) == PieceEnum.PAWN:
                    evaluation += int(PiecePrice.PAWN) * side_sign
                elif PieceEnum(num_piece) == PieceEnum.KNIGHT:
                    evaluation += int(PiecePrice.KNIGHT) * side_sign
                elif PieceEnum(num_piece) == PieceEnum.BISHOP:
                    evaluation += int(PiecePrice.BISHOP) * side_sign
                elif PieceEnum(num_piece) == PieceEnum.ROOK:
                    evaluation += int(PiecePrice.ROOK) * side_sign
                elif PieceEnum(num_piece) == PieceEnum.QUEEN:
                    evaluation += int(PiecePrice.QUEEN) * side_sign
                elif PieceEnum(num_piece) == PieceEnum.KING:
                    evaluation += int(PiecePrice.KING) * side_sign

    return evaluation


def evaluate_piece_placement(position, side=PieceColor.WHITE):
    for num_color in range(0, 9, 8):
        for num_piece in range(1, 7):
            separate_bit_pieces = get_separate_piece(position, PieceEnum(num_piece), PieceColor(num_color))
            separate_num_pieces = [get_num_from_bitboard(num) for num in separate_bit_pieces]
            print(PieceEnum(num_piece), PieceColor(num_color), separate_num_pieces)

            side_sign = 1 if PieceColor(num_color) == side else -1
            # for sep_num_piece in separate_num_pieces:
            for _ in separate_bit_pieces:
                pass  # TODO: maybe the part above should be written only once, and all the evaluation moments in if-s


eval = evaluate_position(position, side=PieceColor.WHITE)
regular_eval = eval / 100

print(regular_eval)
