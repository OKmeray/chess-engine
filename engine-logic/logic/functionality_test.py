from logic.helping_function import get_position_from_fen, GenerateMove, PieceEnum, PieceColor, get_separate_piece, \
    get_all_pieces_of_one_color, get_nums_from_bit_nums, get_square_name_by_num, get_num_from_bitboard

# fen = "3k4/8/8/5N2/8/8/3N4/7K w - - 0 1"
# fen = "3k4/8/8/8/8/8/3N4/7K w - - 0 1"
fen = "3k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"

position = get_position_from_fen(fen)


# without checking if it's legal or not
def get_all_moves(position):
    result = []

    for bit_color in range(0, 9, 8):
        own_bitboard = position.get_white_bitboard() if bit_color else position.get_black_bitboard()
        enemy_bitboard = position.get_white_bitboard() if not bit_color else position.get_black_bitboard()

        for bit_piece in range(1, 7):
            separate_pieces = get_separate_piece(position, PieceEnum(bit_piece), PieceColor(bit_color))
            separate_pieces_moves = []

            for piece in separate_pieces:
                if bit_piece == PieceEnum.PAWN:
                    pass
                    # TODO:
                    # separate_pieces_moves.append(GenerateMove.generate_pawn_move(piece, own_bitboard, enemy_bitboard))
                elif bit_piece == PieceEnum.KNIGHT:
                    separate_pieces_moves.append(GenerateMove.generate_knight_move(piece, own_bitboard))
                elif bit_piece == PieceEnum.BISHOP:
                    separate_pieces_moves.append(GenerateMove.generate_bishop_move(piece, own_bitboard, enemy_bitboard))
                elif bit_piece == PieceEnum.ROOK:
                    separate_pieces_moves.append(GenerateMove.generate_rook_move(piece, own_bitboard, enemy_bitboard))
                elif bit_piece == PieceEnum.QUEEN:
                    separate_pieces_moves.append(GenerateMove.generate_queen_move(piece, own_bitboard, enemy_bitboard))
                elif bit_piece == PieceEnum.KING:
                    separate_pieces_moves.append(GenerateMove.generate_king_move(piece, own_bitboard, enemy_bitboard))

                for piece_moves in separate_pieces_moves:
                    piece_moves_num = get_nums_from_bit_nums(piece_moves)
                    local_dictionary = {bit_color + bit_piece: {"square": get_num_from_bitboard(piece), "possible_moves": piece_moves_num}}
                    result.append(local_dictionary)

                separate_pieces_moves.clear()
    return result


result = get_all_moves(position=position)

for i in result:
    print(i)
