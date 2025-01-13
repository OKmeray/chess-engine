import time
import random
from core.domain.engine import Position
from core.domain.engine.enums import PieceEnum, PieceColor
from core.domain.engine.square_helping_functions import get_square_name_by_num
from core.domain.evaluation.evaluation import evaluate_position
from core.domain.fen.MoveParser import MoveParser
from core.domain.fen.fen_generator import get_position_from_fen, get_fen_from_position
from core.domain.minimax.minimax import minimax


def count_positions(position: Position, depth=2, outcomes=0, is_root=False):
    if position.is_game_over():
        return 1

    if depth == 0:
        return 1

    all_moves = position.get_all_separate_moves()
    for move in all_moves:
        position_copy = position.clone()
        # print(get_fen_from_position(position_copy), move)
        if is_root:
            print(move, end="\t")
        position_copy.apply_move(move)
        position_copy_count = count_positions(position_copy, depth - 1)
        if is_root:
            print(position_copy_count)
        outcomes += position_copy_count

    return outcomes


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


mate_in_1 = '4r2k/1p3rbp/2p1N1p1/p3n3/P2NB1nq/1P6/4R1P1/B1Q2RK1 b - - 4 32'
easy_fen = '5k2/8/5K2/8/2Q5/8/8/8 w - - 0 1'
mate_in_5 = 'R1bq1r1k/6pp/2pP1b2/5P2/4pPB1/1r1nP2K/NP2Q2P/2B3NR b - - 10 24'
mate_in_3 = '4r1k1/4r1p1/8/p2R1P1K/5P1P/1QP3q1/1P6/3R4 b - - 0 1'
mate_in_3_simplified = '6k1/6p1/8/3R1P1K/7P/6q1/8/8 b - - 0 1'
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
rook_and_pawn_vs_queen_ending = '8/1k6/8/3R4/4P3/4K3/q7/8 b - - 0 1'

before_hang_queen = 'rn1qkb1r/pp2nppp/2p1p3/3pP3/3P4/5B1P/PPP2PP1/RNBQ1RK1 b Qkq - 0 1'
after_hang_queen = 'rn2kb1r/pp2nppp/2pqp3/3pP3/3P4/5B1P/PPP2PP1/RNBQ1RK1 w Qkq - 1 2'
# initial_position = get_position_from_fen(before_hang_queen)

# moves = initial_position.get_moves_sorted_by_priority()
# for move in moves:
#     print(move)
# start = time.time()
# best_score, best_move = minimax(initial_position, -float('inf'), float('inf'), depth=3)
# end = time.time()
# from core.domain.minimax.minimax import visited_nodes
# print(visited_nodes)
# print(best_score, best_move)
# print(f"Taken time: {str(round((end - start), 2))} s. ({str(round((end - start) / 60))} m.)")
check_correctness_fen = 'rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8'
starting_position = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'
position_5 = 'rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8'
position_5_qd2 = 'rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPPQNnPP/RNB1K2R b KQ - 2 8'
position_5_qd2_nd3 = 'rnbq1k1r/pp1Pbppp/2p5/8/2B5/3n4/PPPQN1PP/RNB1K2R w KQ - 3 9'
position_5_qd2_nh1 = 'rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPPQN1PP/RNB1K2n w Q - 0 9'
en_passant_legality = "rn2k3/ppppp3/8/8/8/8/8/2K5 b q - 0 1"
check_correctness_position = get_position_from_fen(position_5)
start = time.time_ns()
outcomes = count_positions(check_correctness_position, depth=1, outcomes=0, is_root=True)
end = time.time_ns()
print(outcomes)
print((end - start) / 1e9, "s.")

