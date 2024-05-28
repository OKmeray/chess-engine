from core.domain.engine.enums import PiecePrice, CastleEnum, PieceEnum, PieceColor
from core.domain.engine.square_helping_functions import get_num_from_bitboard


white_pawn_square_table = [
    0,  0,  0,  0,  0,  0,  0,  0,
    50, 50, 50, 50, 50, 50, 50, 50,
    10, 10, 20, 30, 30, 20, 10, 10,
     5,  5, 10, 25, 25, 10,  5,  5,
     0,  0,  0, 20, 20,  0,  0,  0,
     5, -5,-10,  0,  0,-10, -5,  5,
     5, 10, 10,-20,-20, 10, 10,  5,
     0,  0,  0,  0,  0,  0,  0,  0
]
black_pawn_square_table = [i * (-1) for i in white_pawn_square_table[::-1]]

white_knight_square_table = [
    -50,-40,-30,-30,-30,-30,-40,-50,
    -40,-20,  0,  0,  0,  0,-20,-40,
    -30,  0, 10, 15, 15, 10,  0,-30,
    -30,  5, 15, 20, 20, 15,  5,-30,
    -30,  0, 15, 20, 20, 15,  0,-30,
    -30,  5, 10, 15, 15, 10,  5,-30,
    -40,-20,  0,  5,  5,  0,-20,-40,
    -50,-40,-30,-30,-30,-30,-40,-50,
]
black_knight_square_table = [i * (-1) for i in white_knight_square_table[::-1]]


white_bishop_square_table = [
    -20,-10,-10,-10,-10,-10,-10,-20,
    -10,  0,  0,  0,  0,  0,  0,-10,
    -10,  0,  5, 10, 10,  5,  0,-10,
    -10,  5,  5, 10, 10,  5,  5,-10,
    -10,  0, 10, 10, 10, 10,  0,-10,
    -10, 10, 10, 10, 10, 10, 10,-10,
    -10,  5,  0,  0,  0,  0,  5,-10,
    -20,-10,-10,-10,-10,-10,-10,-20,
]
black_bishop_square_table = [i * (-1) for i in white_bishop_square_table[::-1]]


white_rook_square_table = [
      0,  0,  0,  0,  0,  0,  0,  0,
      5, 10, 10, 10, 10, 10, 10,  5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
     -5,  0,  0,  0,  0,  0,  0, -5,
      0,  0,  0,  5,  5,  0,  0,  0
]
black_rook_square_table = [i * (-1) for i in white_rook_square_table[::-1]]


white_queen_square_table = [
    -20,-10,-10, -5, -5,-10,-10,-20,
    -10,  0,  0,  0,  0,  0,  0,-10,
    -10,  0,  5,  5,  5,  5,  0,-10,
     -5,  0,  5,  5,  5,  5,  0, -5,
      0,  0,  5,  5,  5,  5,  0, -5,
    -10,  5,  5,  5,  5,  5,  0,-10,
    -10,  0,  5,  0,  0,  0,  0,-10,
    -20,-10,-10, -5, -5,-10,-10,-20
]
black_queen_square_table = [i * (-1) for i in white_queen_square_table[::-1]]


white_king_square_middlegame_table = [
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -30,-40,-40,-50,-50,-40,-40,-30,
    -20,-30,-30,-40,-40,-30,-30,-20,
    -10,-20,-20,-20,-20,-20,-20,-10,
     20, 20,  0,  0,  0,  0, 20, 20,
     20, 30, 10,  0,  0, 10, 30, 20
]
black_king_square_middlegame_table = [i * (-1) for i in white_king_square_middlegame_table[::-1]]


white_king_square_endgame_table = [
    -50,-40,-30,-20,-20,-30,-40,-50,
    -30,-20,-10,  0,  0,-10,-20,-30,
    -30,-10, 20, 30, 30, 20,-10,-30,
    -30,-10, 30, 40, 40, 30,-10,-30,
    -30,-10, 30, 40, 40, 30,-10,-30,
    -30,-10, 20, 30, 30, 20,-10,-30,
    -30,-30,  0,  0,  0,  0,-30,-30,
    -50,-30,-30,-30,-30,-30,-30,-50
]
black_king_square_endgame_table = [i * (-1) for i in white_king_square_endgame_table[::-1]]


# TODO: check for mobility of pieces (by checking the amount of squares the can move, but if
## if the king is in check it should be calculated in another way)
def evaluate_position(position):
    evaluation = 0

    for num_color in range(0, 9, 8):
        for num_piece in range(1, 7):
            separate_bit_pieces = position.get_separate_piece(PieceEnum(num_piece), PieceColor(num_color))
            separate_pieces_squares = [get_num_from_bitboard(num) for num in separate_bit_pieces]

            # TODO: redo material as
            side_sign = 1 if PieceColor(num_color) == PieceColor.WHITE else -1  # to have evaluation from white side always
            for _ in separate_bit_pieces:
                if PieceEnum(num_piece) == PieceEnum.PAWN:
                    # material
                    evaluation += int(PiecePrice.PAWN) * side_sign

                    # piece placement
                    for square in separate_pieces_squares:
                        evaluation += white_pawn_square_table[square] if side_sign == 1 else black_pawn_square_table[square]

                elif PieceEnum(num_piece) == PieceEnum.KNIGHT:
                    # material
                    evaluation += int(PiecePrice.KNIGHT) * side_sign

                    # piece placement
                    for square in separate_pieces_squares:
                        evaluation += white_knight_square_table[square] if side_sign == 1 else black_knight_square_table[square]

                elif PieceEnum(num_piece) == PieceEnum.BISHOP:
                    # material
                    evaluation += int(PiecePrice.BISHOP) * side_sign

                    # piece placement
                    for square in separate_pieces_squares:
                        evaluation += white_bishop_square_table[square] if side_sign == 1 else black_bishop_square_table[square]

                elif PieceEnum(num_piece) == PieceEnum.ROOK:
                    # material
                    evaluation += int(PiecePrice.ROOK) * side_sign

                    # piece placement
                    for square in separate_pieces_squares:
                        evaluation += white_rook_square_table[square] if side_sign == 1 else black_rook_square_table[square]

                elif PieceEnum(num_piece) == PieceEnum.QUEEN:
                    # material
                    evaluation += int(PiecePrice.QUEEN) * side_sign

                    # piece placement
                    for square in separate_pieces_squares:
                        evaluation += white_queen_square_table[square] if side_sign == 1 else black_queen_square_table[square]

                elif PieceEnum(num_piece) == PieceEnum.KING:
                    # material
                    evaluation += int(PiecePrice.KING) * side_sign

                    # piece placement
                    is_queenless_endgame = (
                        position.get_piece_bitboard(piece=PieceEnum.QUEEN, color=PieceColor.WHITE) |
                        position.get_piece_bitboard(piece=PieceEnum.QUEEN, color=PieceColor.BLACK)
                                   ).bitboard == 0

                    piece_bitboard = (
                        (
                            (
                                    position.get_piece_bitboard(piece=PieceEnum.PAWN, color=PieceColor.WHITE) |
                                    position.get_piece_bitboard(piece=PieceEnum.KING, color=PieceColor.WHITE)
                            ) ^ position.get_white_bitboard()
                        ) |
                        (
                            (
                                    position.get_piece_bitboard(piece=PieceEnum.PAWN, color=PieceColor.BLACK) |
                                    position.get_piece_bitboard(piece=PieceEnum.KING, color=PieceColor.BLACK)
                            ) ^ position.get_black_bitboard()
                        )
                    )

                    count = 0
                    while piece_bitboard.bitboard:
                        count += piece_bitboard.bitboard & 1
                        piece_bitboard.bitboard >>= 1

                    is_queen_endgame = False
                    if count <= 4:
                        is_queen_endgame = True

                    # piece placement
                    if is_queen_endgame or is_queenless_endgame:
                        for square in separate_pieces_squares:
                            evaluation += white_king_square_endgame_table[square] if side_sign == 1 else black_king_square_endgame_table[square]
                    else:
                        for square in separate_pieces_squares:
                            evaluation += white_king_square_middlegame_table[square] if side_sign == 1 else black_king_square_middlegame_table[square]

    return evaluation


# # TODO: calculate piece activity (by checking how many possible moves each piece has)
#
#
# from functionality_test import position
# eval = evaluate_position(position, side=PieceColor.WHITE)
# regular_eval = eval / 100
#
# print(regular_eval)
