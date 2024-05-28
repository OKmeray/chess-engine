# import pytest
# from core.domain.fen import get_position_from_fen
# from data_preparation import test_data_for_mate_in_one
# from logic.mcts.MCTS import MCTS
#

# @pytest.mark.parametrize("fen, best_move_detail", test_data_for_mate_in_one)
# def test_best_move_seek(fen, best_move_detail):
#     position = get_position_from_fen(fen)
#
#     iterations = 30  # TODO: maybe transfer from test_data
#
#     mcts = MCTS(position)
#     mcts.run(iterations)
#     best_move = mcts.best_move()
#
#     assert best_move.move["piece"] == best_move_detail["piece"]
#     assert best_move.move["color"] == best_move_detail["color"]
#     assert best_move.move["square"] == best_move_detail["square"]
#     assert best_move.move["move"] == best_move_detail["move"]

