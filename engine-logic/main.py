import time

from logic.engine.enums import PieceEnum, PieceColor, CastleEnum
from logic.evaluation.evaluation import evaluate_position
from logic.fen.fen_generator import get_position_from_fen
from logic.engine.square_helping_functions import get_square_name_by_num
from logic.evaluation.evaluation import evaluate_position
from logic.fen.fen_generator import get_position_from_fen
from logic.mcts.MCTS import MCTS


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


# fen = "3k4/8/8/5N2/8/8/3N4/7K w - - 0 1"
# fen = "3k4/8/8/8/8/8/3N4/7K w - - 0 1"
# fen = "3k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"
# fen = "3k4/8/3b4/3P1N2/2pr4/4N1B1/8/7K w - - 0 1"  # fen_without_queens
# fen = "q2k4/8/3b4/3P1N2/2pr4/3QN1B1/8/7K w - - 0 1"  # fen with two queens
# new_fen = "rnb1kb1r/2qp1p1p/pP6/3npPp1/4P3/1p6/P3K1PP/RNBQ1BNR w kq g6 0 11"


# black_not_checkmated = "K7/8/7r/8/8/8/1q6/5k2 w - - 0 1"
# black_is_checkmated = "K7/8/r7/8/8/8/1q6/5k2 w - - 0 1"
# regular_position = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
# position = get_position_from_fen(regular_position)
#
# print(position.is_checkmate())
# print(evaluate_position(position, position.side_to_move))  # TODO: correct evaluate_position
# # print(position.is_king_in_check())
# # print(position.get_all_moves())
#
# move_detail = {
#     'piece': PieceEnum.KING,
#     'color': PieceColor.WHITE,
#     'square': 1,
#     'move': 2
# }
# print(position.is_move_legal(move_detail))
# for i in result:
#     print(i)

# UI_moves = return_moves_for_UI(result)
# print(UI_moves)


# print("\n\nRegular notation: \n")
# regular_notation_res = return_moves_in_regular_notation(result)
# for i in regular_notation_res:
#     print(i)


# print(45)
# print(get_num_from_bitboard(get_bitboard_from_num(45)))
#
# def main():
#     print("For the time being empty")
#
#
# fen = "b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/2R2r1K w - - 0 27"
#
# position_to_test = get_position_from_fen(fen)
# # print(position_to_test.is_checkmate())
#
# new_fen = "K3R3/Q5B1/8/8/8/2br4/1q6/5k2 w - - 0 1"
# position_for_apply_move = get_position_from_fen(new_fen)
# for piece in position_for_apply_move.generate_all_pseudo_legal_moves():
#     print(piece)
# move_detail = {'piece': PieceEnum.BISHOP, 'color': PieceColor.WHITE, 'square': 14, 'move': 42}
# position_for_apply_move.apply_move(move_detail)  # TODO: correct apply_move
#
# print("\n")
# for piece in position_for_apply_move.generate_all_pseudo_legal_moves():
#     print(piece)

# fen = "r1bqk2r/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7"
#
# position_to_test = get_position_from_fen(fen)
# move = {
#              'piece': PieceEnum.KING,
#              'color': PieceColor.BLACK,
#              'square': 4,
#              'move': 6
#         }
# print(position_to_test.is_move_legal(move))
# print(position_to_test.generate_all_pseudo_legal_moves())
# # print(position_to_test.get_all_moves())

# fen_init = "r1bqk2r/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7"
# fen_end = "r1bq1rk1/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 w - - 2 8"
# move_detail = {
#      'piece': PieceEnum.KING,
#      'color': PieceColor.BLACK,
#      'square': 4,
#      'move': 6
# }
# position = get_position_from_fen(fen_init)
# print(position.castling_rights[CastleEnum.BlackShortCastle])
# position.apply_move(move_detail)
# print(position.castling_rights[CastleEnum.BlackShortCastle])
# print(position.get_all_bitboard().get_8by8_board())
# position_end = get_position_from_fen(fen_end)
# print()
# print(position_end.get_all_bitboard().get_8by8_board())
#
# fen_init = "8/8/3k4/8/4N3/3N4/3K4/8 w - - 0 1"
# move_detail = {
#              'piece': PieceEnum.KING,
#              'color': PieceColor.WHITE,
#              'square': 60,
#              'move': 62
#         }
# position = get_position_from_fen(fen_init)
# print(position.is_king_in_check())
# # print(position.is_move_legal(move_detail))
# #
# # print(position.get_all_bitboard().bitboard.bit_count())
#

# initial_position = Position()
# ERROR_FEN = "4r2k/1Bq2rbp/6p1/p1p1n3/P2p2n1/1P1N1N2/6P1/B1Q1RRK1 b - - 0 31"
fen = "8/8/8/5q2/8/2k5/8/2K5 b - - 10 24"
initial_position = get_position_from_fen(fen)
mcts = MCTS(initial_position)
print("RUN")
start = time.time()
mcts.run(1500)
end = time.time()
print(end - start, "seconds")
print("RUN")
# best_move_node = mcts.best_move()
# print(best_move_node.move)
print(len(mcts.root.children))
for child in mcts.root.children:
    print(child.move, child.wins, "/", child.visits, "\t[", ["".join([str(i.wins), "/", str(i.visits)]) for i in child.children], "]")
