import os
import xml.etree.ElementTree as ET

from logic.engine.enums import PieceEnum, PieceColor
from logic.fen.generator import get_position_from_fen


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
# test_data_for_is_threefold_repetition = []
# test_data_for_is_insufficient_material = []

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
    ("K3R3/Q5B1/8/8/8/2br4/1q6/5k2 w - - 0 1", "K3R3/Q7/8/8/8/2Br4/1q6/5k2 b - - 1 1", {'piece': PieceEnum.BISHOP, 'color': PieceColor.WHITE, 'square': 14, 'move': 42}),
    ("K3R3/Q5B1/8/8/8/2br4/1q6/5k2 w - - 0 1", "K3R3/5QB1/8/8/8/2br4/1q6/5k2 b - - 1 1", {'piece': PieceEnum.QUEEN, 'color': PieceColor.WHITE, 'square': 8, 'move': 13}),
    ("b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/2R2r1K w - - 0 27", "b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K b - - 0 27", {'piece': PieceEnum.ROOK, 'color': PieceColor.WHITE, 'square': 58, 'move': 61}),
    ("b1r4k/2q3pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K b - - 0 27", "bqr4k/6pp/pp2Q3/3N4/1PPbB3/P7/6PP/5R1K w - - 1 28", {'piece': PieceEnum.QUEEN, 'color': PieceColor.BLACK, 'square': 10, 'move': 1}),
]

