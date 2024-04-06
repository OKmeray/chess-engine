import os
import xml.etree.ElementTree as ET
from logic.engine.enums import PieceEnum, PieceColor


def parse_test_data_from_xml(xml_file_path):
    dir_path = os.path.dirname(os.path.realpath(__file__))
    xml_file_path = os.path.join(dir_path, xml_file_path)

    tree = ET.parse(xml_file_path)
    root = tree.getroot()

    data = []
    for position in root.findall('position'):
        fen = position.find('fen').text
        state = position.find('state').text
        result = position.find('result').text
        data.append({"fen": fen, "state": state, "result": result})

    return data


xml_file_path = "test_data.xml"
game_collection = parse_test_data_from_xml(xml_file_path)

test_data_for_is_king_in_check = []
test_data_for_is_checkmate = []

test_data_for_is_stalemate = []
test_data_for_is_fifty_move_rule_reached = []
test_data_for_is_insufficient_material = []
# test_data_for_is_threefold_repetition = []


# checks
for game in game_collection:
    if game["state"] == "check" or game["state"] == "checkmate":
        test_data_for_is_king_in_check.append((game["fen"], True))
    else:
        test_data_for_is_king_in_check.append((game["fen"], False))

# checkmates
for game in game_collection:
    if game["state"] == "checkmate":
        test_data_for_is_checkmate.append((game["fen"], True))
    else:
        test_data_for_is_checkmate.append((game["fen"], False))

# stalemates
for game in game_collection:
    if game["state"] == "stalemate":
        test_data_for_is_stalemate.append((game["fen"], True))
    else:
        test_data_for_is_stalemate.append((game["fen"], False))

# fifty_move_rule_reached
for game in game_collection:
    if game["state"] == "draw by insufficient material":
        test_data_for_is_insufficient_material.append((game["fen"], True))
    else:
        test_data_for_is_insufficient_material.append((game["fen"], False))

# insufficient_material
for game in game_collection:
    if game["state"] == "draw by 50-move rule":
        test_data_for_is_fifty_move_rule_reached.append((game["fen"], True))
    else:
        test_data_for_is_fifty_move_rule_reached.append((game["fen"], False))

# #
# for game in game_collection:
#     if game["state"] == "":
#         .append((game["fen"], True))
#     else:
#         .append((game["fen"], False))
#
# #
# for game in game_collection:
#     if game["state"] == "":
#         .append((game["fen"], True))
#     else:
#         .append((game["fen"], False))

test_data_for_apply_move = [
    (
        "K3R3/Q5B1/8/8/8/2br4/1q6/5k2 w - - 0 1",
        "K3R3/Q7/8/8/8/2Br4/1q6/5k2 b - - 0 1",
        {
            'piece': PieceEnum.BISHOP,
            'color': PieceColor.WHITE,
            'square': 14,
            'move': 42
        }
    ),

    (
        "K3R3/Q5B1/8/8/8/2br4/1q6/5k2 w - - 0 1",
        "K3R3/5QB1/8/8/8/2br4/1q6/5k2 b - - 1 1",
         {
             'piece': PieceEnum.QUEEN,
             'color': PieceColor.WHITE,
             'square': 8,
             'move': 13
         }
    ),

    (
        "b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/2R2r1K w - - 0 27",
        "b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K b - - 0 27",
        {
            'piece': PieceEnum.ROOK,
            'color': PieceColor.WHITE,
            'square': 58,
            'move': 61
        }
    ),

    (
        "b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K b - - 0 27",
        "bqr4k/6pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K w - - 1 28",
        {
            'piece': PieceEnum.QUEEN,
            'color': PieceColor.BLACK,
            'square': 10,
            'move': 1
        }
    ),

    (
        "r1b1k2r/4bppp/p1pq1n2/1pp1p3/4P2B/3P1N1P/PPP2PP1/RN1QK2R b KQkq - 6 10",
        "r1b1k2r/4bppp/p1pq1n2/1p2p3/2p1P2B/3P1N1P/PPP2PP1/RN1QK2R w KQkq - 0 11",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.BLACK,
             'square': 26,
             'move': 34
        }
    ),

    (
        "r1b1k2r/4bppp/p1p2n2/1p2p3/1qp1P3/3P1NBP/PPPN1PP1/R2QK2R b KQkq - 3 12",
        "r1b1k2r/4bppp/p1p2n2/1p2p3/2p1P3/3P1NBP/PqPN1PP1/R2QK2R w KQkq - 0 13",
        {
             'piece': PieceEnum.QUEEN,
             'color': PieceColor.BLACK,
             'square': 33,
             'move': 49
        }
    ),

    (
        "r1b1k2r/4bppp/p1p2n2/1p2B3/2p1P3/3P1N1P/PqPN1PP1/R2QK2R b KQkq - 0 13",
        "r1b1k2r/4bppp/p1p2n2/1p2B3/4P3/3p1N1P/PqPN1PP1/R2QK2R w KQkq - 0 14",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.BLACK,
             'square': 34,
             'move': 43
        }
    ),

    # applying castles
    (
        "r1bqk2r/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7",
        "r1bq1rk1/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 w - - 2 8",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.BLACK,
             'square': 4,
             'move': 6
        }
    ),

    (
        "r3kb1r/pppppp1p/8/8/8/1QPq4/PP1PPP1P/R3K2R w KQkq - 0 1",
        "r3kb1r/pppppp1p/8/8/8/1QPq4/PP1PPP1P/2KR3R b kq - 1 1",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.WHITE,
             'square': 60,
             'move': 58
        }
    ),

    # applying en passant
    (
        "r1bqkbnr/ppppp1pp/2n5/4Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 2",
        "r1bqkbnr/ppppp1pp/2n2P2/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.WHITE,
             'square': 28,
             'move': 21
        }
    ),

    (
        "r1bqkbnr/ppp1p2p/2n2p2/8/2Pp4/5N2/PP1P1PPP/RNBQKB1R b KQkq c3 0 2",
        "r1bqkbnr/ppp1p2p/2n2p2/8/8/2p2N2/PP1P1PPP/RNBQKB1R w KQkq - 0 3",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.BLACK,
             'square': 35,
             'move': 42
        }
    ),

    # applying promotion
    (
        "3r4/2P2q2/6k1/8/8/8/3K4/8 w - - 0 1",
        "3Q4/5q2/6k1/8/8/8/3K4/8 b - - 0 1",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.WHITE,
             'square': 10,
             'move': 3
        }
    ),

    (
        "5nrq/2P2p2/1P4k1/6p1/8/B3Q3/3K4/8 w - - 0 1",
        "2Q2nrq/5p2/1P4k1/6p1/8/B3Q3/3K4/8 b - - 0 1",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.WHITE,
             'square': 10,
             'move': 2
        }
    ),

    (
        "5nr1/K1P1q2k/1PQ5/6p1/8/B7/5p2/8 b - - 0 1",
        "5nr1/K1P1q2k/1PQ5/6p1/8/B7/8/5q2 w - - 0 2",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.BLACK,
             'square': 53,
             'move': 61
        }
    ),

    (
        "5nr1/K1P1q2k/1PQ5/6p1/8/B7/5p2/4N3 b - - 0 2",
        "5nr1/K1P1q2k/1PQ5/6p1/8/B7/8/4q3 w - - 0 3",
        {
             'piece': PieceEnum.PAWN,
             'color': PieceColor.BLACK,
             'square': 53,
             'move': 60
        }
    ),
]

# castles check
test_data_for_is_move_legal = [
    (
        "rnbqkb1r/pppppppp/8/8/8/4B3/PPPPPPPP/RNBQK1nR w KQkq - 0 1",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.WHITE,
             'square': 60,
             'move': 62
        },
        False
    ),

    (
        "r1bqk2r/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.BLACK,
             'square': 4,
             'move': 6
        },
        True
    ),

    (
        "r2qkb1r/1pp2ppp/p1np1n2/4p3/B1b1P3/P4N2/1PPP1PPP/RNBQK2R w Kkq - 4 8",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.WHITE,
             'square': 60,
             'move': 62
        },
        False
    ),

    (
        "r1bq1rk1/p1pp1ppp/1pn2n2/1Bb1p3/4PP2/3P4/PPPN2PP/RNBQK2R w KQ - 1 7",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.WHITE,
             'square': 60,
             'move': 62
        },
        False
    ),

    (
        "rnbq1knr/ppppp1pp/8/8/8/8/PPPPP1PP/RNBQK2R w KQkq - 0 1",
        {
             'piece': PieceEnum.KING,
             'color': PieceColor.WHITE,
             'square': 60,
             'move': 62
        },
        True
    ),

    (
        "rnbqkbr1/pppppp1p/8/8/8/8/PPPPPP1P/RNBQK2R w KQq - 0 1",
        {
            'piece': PieceEnum.KING,
            'color': PieceColor.WHITE,
            'square': 60,
            'move': 62
        },
        False
    ),

    (
        "rn1qkb1r/pppppp1p/8/8/8/7b/PPPPPP1P/RNBQK2R w KQkq - 0 1",
        {
            'piece': PieceEnum.KING,
            'color': PieceColor.WHITE,
            'square': 60,
            'move': 62
        },
        False
    ),

    (
        "rn2kb2/pppppp1p/8/8/8/4bq1r/PPPPPP1P/RNBQK2R w KQhq - 0 1",
        {
            'piece': PieceEnum.KING,
            'color': PieceColor.WHITE,
            'square': 60,
            'move': 62
        },
        True
    ),

    (
        "r3kb2/pppppp1p/2n5/8/8/4bq1r/PPPPPP1P/RNBQK2R b KQq - 0 1",
        {
            'piece': PieceEnum.KING,
            'color': PieceColor.BLACK,
            'square': 4,
            'move': 2
        },
        True
    ),

    (
        "r3kb1r/pppppp1p/8/8/8/1QPq4/PP1PPP1P/R3K2R w KQkq - 0 1",
        {
            'piece': PieceEnum.KING,
            'color': PieceColor.WHITE,
            'square': 60,
            'move': 58
        },
        True
    )

]

print(test_data_for_is_king_in_check)
