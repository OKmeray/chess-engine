import time

# from logic.engine.enums import PieceEnum, PieceColor, CastleEnum
import random

from core.domain.engine.enums import PieceEnum, PieceColor
from core.domain.engine.square_helping_functions import get_square_name_by_num
from core.domain.evaluation.evaluation import evaluate_position
from core.domain.fen.MoveParser import MoveParser
from core.domain.fen.fen_generator import get_position_from_fen, get_fen_from_position
from core.domain.minimax.minimax import minimax


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


# fen = "8/8/8/5q2/8/2k5/8/2K5 b - - 10 24"
# initial_position = get_position_from_fen(fen)
# mcts = MCTS(initial_position)
#
# start = time.time()
# mcts.run(15)
# end = time.time()
# print(end - start, "seconds")
#
# print(len(mcts.root.children))
# print(mcts.root.visits, mcts.root.wins)
# for child in mcts.root.children:
#     print(child.move, child.wins, "/", child.visits, "\t[", ["".join([str(i.wins), "/", str(i.visits)]) for i in child.children], "]")
# print("BEST MOVE IS:", mcts.best_move().move)

mate_in_1 = '4r2k/1p3rbp/2p1N1p1/p3n3/P2NB1nq/1P6/4R1P1/B1Q2RK1 b - - 4 32'
easy_fen = '5k2/8/5K2/8/2Q5/8/8/8 w - - 0 1'
mate_in_5 = 'R1bq1r1k/6pp/2pP1b2/5P2/4pPB1/1r1nP2K/NP2Q2P/2B3NR b - - 10 24'
mate_in_3 = '4r1k1/4r1p1/8/p2R1P1K/5P1P/1QP3q1/1P6/3R4 b - - 0 1'
mate_in_3_simplified = '6k1/6p1/8/3R1P1K/7P/6q1/8/8 b - - 0 1'

rook_and_pawn_vs_queen_ending = '8/1k6/8/3R4/4P3/4K3/q7/8 b - - 0 1'

mate_in_4 = 'rn4k1/p6q/bp2p2p/3pPN2/P3n2Q/1Pb1RN1P/5PBK/3R4 w - - 1 27'
mate_in_4_2 = '4r1rk/1bq5/pp5p/P2pPp2/3Nn2Q/1P1n1N1P/2R2P1K/5BR1 w - - 2 34'
mate_in_4_3 = '2r1k3/pp1b3r/1qnPp2p/1Pb3pP/6Pn/PR1Q2N1/2P1B3/2B2R1K w - - 1 29'
mate_in_4_4 = '2r3k1/q2P1r1p/2b1Nb2/1p1R1p2/pBn5/6P1/6BK/3QR3 w - - 1 37'
mate_in_8_1 = '4r2k/1p3rbp/2p1N1pn/p3n3/P2NB3/1P4q1/4R1P1/B1Q2RK1 b - - 4 32'
mate_in_8_2 = 'rn3rk1/p7/bpq1pp1Q/3pP3/P2Nn3/1Pb2NPP/5PB1/3RR1K1 w - - 1 22'
mate_in_8_3 = '1Rbq1r1k/6pp/2pP4/2n2P2/2Q2P1b/2r1Pp2/1P5P/2B2KNR b - - 1 22'
mate_in_8_4 = 'R1b4k/6pp/2pP4/1rn2rq1/N1Q2P1b/4PN2/1P5P/2B2K1R w - - 1 24'
mate_in_8_5 = 'rn4k1/5rbn/1p1p4/1p1q1p1p/3P4/P1B1P2P/Q2N1PR1/1K4R1 w - - 2 27'
mate_in_8_6 = '4r1rk/p6q/bp2pP1Q/3pP3/P2Nn3/1Pbn1N1P/5PBK/2R3R1 b - - 0 30'
mate_in_8_7 = 'r4k2/pb6/1p2qP1p/P2p3r/5Q2/1Pn1RN1P/5PB1/6K1 w - - 1 33'
mate_in_8_26 = '2r1k2r/1p1b4/pq2p2p/1Pb3pP/2P1NnPn/1R1Q4/P3B3/2BR3K w - - 2 30'
mate_in_8_27 = '3q3k/1p6/p2N1n1p/1Pr1pBpb/PBPb4/5Q2/8/1R5K b - - 3 42'
print(mate_in_8_27)
new_pos = '2rnk3/pp1b3r/1q1Pp2p/1Pb3pP/6Pn/PRQ3N1/1BP1B3/5R1K b - - 4 30'
initial_position = get_position_from_fen(mate_in_8_27)


# double_prev_moves = initial_position.generate_all_pseudo_legal_moves()
# for move in double_prev_moves:
#     print(move)
#
# print("\n")
#
# prev_moves = initial_position.get_all_moves()
# for move in prev_moves:
#     print(move)
#
# print("\n")
#
# moves = initial_position.get_all_separate_moves()
# for move in moves:
#     print(move)

# moves_sorted_by_priority = initial_position.get_moves_sorted_by_priority()
# for move in moves_sorted_by_priority:
#     print(move)
#
# print()
# moves = initial_position.get_moves_sorted_by_priority()
# moves = initial_position.get_all_separate_moves()
# for move in moves:
#     print(move)
# print()
#
# start = time.time()
# best_score, best_move = minimax(initial_position, -float('inf'), float('inf'), depth=5)
# end = time.time()
# print(best_score, best_move)
# print(f"Taken time: {str(round((end - start), 2))} s. ({str(round((end - start) / 60))} m.)")


import re

san_pattern = re.compile(r'^(?P<piece>[KQRBN])?(?P<from_file>[a-h])?(?P<from_rank>[1-8])?(?P<capture>x)?(?P<to_file>[a-h])(?P<to_rank>[1-8])(?P<promotion>=Q|=R|=B|=N)?(?P<check>[+#])?$')
castling_pattern = re.compile(r'^(O-O|O-O-O)$')

# Example mapping for piece characters
piece_map = {
    'P': PieceEnum.PAWN,
    'N': PieceEnum.KNIGHT,
    'B': PieceEnum.BISHOP,
    'R': PieceEnum.ROOK,
    'Q': PieceEnum.QUEEN,
    'K': PieceEnum.KING
}

# Convert file and rank to board index
def convert_file_rank_to_square(file, rank):
    file_index = ord(file) - ord('a')
    rank_index = int(rank) - 1
    return rank_index * 8 + file_index


# def parse_san(san):
#     if castling_pattern.match(san):
#         return {'castling': san}
#     match = san_pattern.match(san)
#     if not match:
#         raise ValueError(f"Invalid SAN move: {san}")
#     return match.groupdict()
# fen = 'rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1'
# move1 = {'fen': 'rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1', 'move': 'Qe7e5'}
# move2 = {'fen': 'rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1', 'move': 'e7e5'}
# move_parser = MoveParser()
# move_detail = move_parser.convert_san_to_move_detail(fen, move2['move'])
# fen2 = "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
# position = get_position_from_fen(fen2)
# fen3 = get_fen_from_position(position)
# print(fen2 == fen3)
# print(fen2)
# print(move_detail)


# fen = '8/8/3k1p2/4P3/3K4/8/8/8 b - - 0 1'
# position = get_position_from_fen(fen)
# all_moves = position.get_all_separate_moves()
# for move in all_moves:
#     print(move)
# TODO: best_move is from the end of the game not from the start of the position
# TODO: order the moves (checks and captures first)
#
#
# position1 = get_position_from_fen("rnbqkbn1/pppppppr/8/8/3PP1p1/8/PPP2PPP/RNB1KBNR w KQq - 0 4")
# position2 = get_position_from_fen("rnbqkbn1/ppppppp1/7r/7p/3PP1Q1/8/PPP2PPP/RNB1KBNR w KQq - 3 4")
# position_starting = get_position_from_fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
# print(evaluate_position(position1))
# print(evaluate_position(position2))
# print(evaluate_position(position_starting))


chess_moves_fen = "rnbqkbnr/ppppppp1/7p/8/P7/5N2/1PPPPPPP/RNBQKB1R b KQkq - 1 2"
chess_moves_fen = "rnbqkbnr/p1ppppp1/7p/1p6/P7/5NP1/1PPPPP1P/RNBQKB1R b KQkq - 0 3"
chess_moves_fen = "rnb1kbnr/2qpppp1/2p4p/pp5P/P7/5NP1/1PPPPP2/RNBQKB1R b KQkq - 0 6"
chess_moves_fen = "rn2kbnr/1bqpppp1/2p4p/pp5P/PP6/5NP1/2PPPP2/RNBQKB1R b KQkq - 0 7"
chess_moves_fen = "rn2kbnr/1b2p1p1/1qpp1p1p/pp5P/PP5N/5PP1/2PPP3/RNBQKB1R w KQkq - 0 10"
chess_moves_fen = "rn2kbnr/1b4p1/1qpppp1p/pp5P/PP5N/3P1PP1/2P1P3/RNBQKB1R w KQkq - 0 11"
chess_moves_fen = "rn2kbnr/1b6/1qpppp1p/pp4pP/PP5N/2PP1PP1/4P3/RNBQKB1R w KQkq g6 0 12"
# chess_moves_fen = "rn3bnr/1b1k4/1qpppp1p/pp4pP/PP3P1N/2PP2P1/4P3/RNBQKB1R w KQ - 1 13"
position = get_position_from_fen(chess_moves_fen)
print(position.en_passant_square)
position.apply_move({'piece': PieceEnum.ROOK, 'color': PieceColor.WHITE, 'square': 56, 'move': 48})
position.apply_move({'piece': PieceEnum.KING, 'color': PieceColor.BLACK, 'square': 4, 'move': 3})
print(position.en_passant_square)
moves = position.get_all_moves()

for move in moves:
    print(move)
