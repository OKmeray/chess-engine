import pytest
from logic.fen.fen_generator import get_position_from_fen
from data_preparation import test_data_for_is_king_in_check, test_data_for_is_checkmate, \
    test_data_for_is_fifty_move_rule_reached, test_data_for_is_stalemate, test_data_for_apply_move, \
    test_data_for_is_insufficient_material


# TODO: tests is_legal_move, legal_moves, pseudo_moves
@pytest.mark.parametrize("fen, expected", test_data_for_is_king_in_check)
def test_is_king_in_check(fen, expected):
    position = get_position_from_fen(fen)
    assert position.is_king_in_check() == expected
    assert True is True





@pytest.mark.parametrize("fen, expected", test_data_for_is_checkmate)
def test_is_checkmate(fen, expected):
    position = get_position_from_fen(fen)
    assert position.is_checkmate() == expected


@pytest.mark.parametrize("fen, expected", test_data_for_is_stalemate)
def test_is_stalemate(fen, expected):
    position = get_position_from_fen(fen)
    assert position.is_stalemate() == expected


@pytest.mark.parametrize("fen, expected", test_data_for_is_fifty_move_rule_reached)
def test_is_fifty_move_rule_reached(fen, expected):
    position = get_position_from_fen(fen)
    assert position.is_fifty_move_rule_reached() == expected


@pytest.mark.parametrize("initial_fen, final_fen, move_detail", test_data_for_apply_move)
def test_apply_move(initial_fen, final_fen, move_detail):
    init_position = get_position_from_fen(initial_fen)
    final_position = get_position_from_fen(final_fen)

    init_position.apply_move(move_detail)

    assert init_position.get_white_bitboard().bitboard == final_position.get_white_bitboard().bitboard
    assert init_position.get_black_bitboard().bitboard == final_position.get_black_bitboard().bitboard
    assert init_position.get_all_bitboard().bitboard == final_position.get_all_bitboard().bitboard
    assert init_position.half_moves == final_position.half_moves
    assert init_position.current_turn == final_position.current_turn


@pytest.mark.parametrize("fen, expected", test_data_for_is_insufficient_material)
def test_is_insufficient_material(fen, expected):
    position = get_position_from_fen(fen)
    assert position.is_insufficient_material() == expected

