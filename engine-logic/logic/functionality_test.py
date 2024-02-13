from logic.helping_function import get_position_from_fen, GenerateMove, PieceEnum, PieceColor, get_separate_piece, \
    get_all_pieces_of_one_color, get_nums_from_bit_nums, get_square_name_by_num

#fen = "3k4/8/8/5N2/8/8/3N4/7K w - - 0 1"
#fen = "3k4/8/8/8/8/8/3N4/7K w - - 0 1"
fen = "3k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"

position = get_position_from_fen(fen)
arr = get_all_pieces_of_one_color(position, PieceColor.WHITE)

separate_knights = get_separate_piece(position, PieceEnum.KNIGHT, PieceColor.WHITE)

all_knights_moves = []
own_bitboard = position.get_white_bitboard()
for knight in separate_knights:
    all_knights_moves.append(GenerateMove.generate_knight_move(own_bitboard, knight))

print(all_knights_moves)

for knight_moves in all_knights_moves:
    knight_moves_num = get_nums_from_bit_nums(knight_moves)
    print([get_square_name_by_num(move) for move in knight_moves_num])

