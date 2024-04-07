import pytest
from logic.fen.fen_generator import get_position_from_fen
from data_preparation import test_data_for_is_move_legal


# TODO: tests is_legal_move, legal_moves, pseudo_moves

# @pytest.mark.parametrize("fen, move_detail, expected", test_data_for_is_move_legal)
# def test_best_move_seek(fen, move_detail, expected):
#     position = get_position_from_fen(fen)
#     assert position.is_move_legal(move_detail) == expected