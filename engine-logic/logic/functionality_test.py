from logic.square_helping_function import  \
   get_nums_from_bit_nums, get_square_name_by_num, get_num_from_bitboard
from logic.position_helping_function import get_position_from_fen, GenerateMove, PieceEnum, PieceColor, get_separate_piece, get_all_pieces_of_one_color

# fen = "3k4/8/8/5N2/8/8/3N4/7K w - - 0 1"
# fen = "3k4/8/8/8/8/8/3N4/7K w - - 0 1"
fen = "3k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"
# fen = "3k4/8/3b4/3P1N2/2pr4/4N1B1/8/7K w - - 0 1"  # fen_without_queens
# fen = "q2k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"  # fen with two queens
new_fen = "rnb1kb1r/2qp1p1p/pP6/3npPp1/4P3/1p6/P3K1PP/RNBQ1BNR w kq g6 0 11"
position = get_position_from_fen(new_fen)


# without checking if it's legal or not
def get_all_moves_by_piece(position):
    result = []

    # for num_color in range(0, 9, 8):
    num_color = position.side_to_move
    own_bitboard = position.get_white_bitboard() if num_color else position.get_black_bitboard()
    enemy_bitboard = position.get_white_bitboard() if not num_color else position.get_black_bitboard()

    for num_piece in range(1, 7):
        separate_bit_pieces = get_separate_piece(position, PieceEnum(num_piece), PieceColor(num_color))
        separate_pieces_moves = []

        for bit_piece in separate_bit_pieces:
            if num_piece == PieceEnum.PAWN:
                bitboard_en_passant_square = None
                if position.en_passant_square is not None:
                    bitboard_en_passant_square = 2 ** position.en_passant_square
                separate_pieces_moves.append(GenerateMove.generate_pawn_move(bit_piece, own_bitboard, enemy_bitboard, color=num_color, en_passant_square=bitboard_en_passant_square))
            elif num_piece == PieceEnum.KNIGHT:
                separate_pieces_moves.append(GenerateMove.generate_knight_move(bit_piece, own_bitboard))
            elif num_piece == PieceEnum.BISHOP:
                separate_pieces_moves.append(GenerateMove.generate_bishop_move(bit_piece, own_bitboard, enemy_bitboard))
            elif num_piece == PieceEnum.ROOK:
                separate_pieces_moves.append(GenerateMove.generate_rook_move(bit_piece, own_bitboard, enemy_bitboard))
            elif num_piece == PieceEnum.QUEEN:
                separate_pieces_moves.append(GenerateMove.generate_queen_move(bit_piece, own_bitboard, enemy_bitboard))
            elif num_piece == PieceEnum.KING:
                separate_pieces_moves.append(GenerateMove.generate_king_move(bit_piece, own_bitboard, enemy_bitboard))

            # print(separate_pieces_moves)
            for piece_moves in separate_pieces_moves:
                # if piece_moves is None:
                piece_moves_num = get_nums_from_bit_nums(piece_moves)
                local_dictionary = {
                    "piece": PieceEnum(num_piece),
                    "color": PieceColor(num_color),
                    "square": get_num_from_bitboard(bit_piece),
                    "possible_moves": piece_moves_num
                }
                result.append(local_dictionary)

            separate_pieces_moves.clear()
    return result


result = get_all_moves_by_piece(position=position)

# for i in result:
#     print(i)


def return_moves_for_UI(all_moves_by_piece):
    result = []

    for piece in all_moves_by_piece:
        piece_letter = None

        match piece['piece']:
            case 1:
                piece_letter = 'p'
            case 2:
                piece_letter = 'n'
            case 3:
                piece_letter = 'b'
            case 4:
                piece_letter = 'r'
            case 5:
                piece_letter = 'q'
            case 6:
                piece_letter = 'k'

        if piece['color'] == 8:
            piece_letter = piece_letter.upper()

        result.append({'piece': piece_letter, 'square': piece['square'], 'possible_moves': piece['possible_moves']})

    return result

# UI_moves = return_moves_for_UI(result)
# print(UI_moves)


def return_moves_in_regular_notation(all_moves_by_piece):
    result = []

    for piece in all_moves_by_piece:
        piece_letter = None

        match piece['piece']:
            case 1:
                piece_letter = 'pawn'
            case 2:
                piece_letter = 'knight'
            case 3:
                piece_letter = 'bishop'
            case 4:
                piece_letter = 'rook'
            case 5:
                piece_letter = 'queen'
            case 6:
                piece_letter = 'king'

        if piece['color'] == 8:
            piece_letter = piece_letter.upper()

        moves_notation = [get_square_name_by_num(piece) for piece in piece['possible_moves']]

        result.append({'piece': piece_letter, 'square': get_square_name_by_num(piece['square']), 'possible_moves': moves_notation})

    return result

print("\n\nRegular notation: \n")
regular_notation_res = return_moves_in_regular_notation(result)
for i in regular_notation_res:
    print(i)