# from core.application.repositories.unit_of_work import IUnitOfWork
# from core.application.schemas.game import GameCreateSchema
# from core.domain.fen.fen_generator import get_position_from_fen, get_fen_from_position
# from core.domain.minimax.minimax import minimax
#
#
# class GameService:
#     def __init__(self, uow: IUnitOfWork):
#         self.uow = uow
#
#     def create_game(self, game: GameCreateSchema):
#         with self.uow as uow:
#             game_dict = {
#                 "created_on": game.created_on,
#                 "user_id": game.user_id
#             }
#             game_db = uow.game.add_one(game_dict)
#             uow.commit()
#             return game_db
#
#     def is_user_move_legal(self, fen, from_square, to_square):
#         position = get_position_from_fen(fen)
#
#         piece, color = position.get_piece_and_color_by_square(int(from_square))
#
#         move = {'piece': piece, 'color': color, 'square': int(from_square), 'move': int(to_square)}
#         if position.is_move_legal(move):
#             position.apply_move(move)
#             return True, get_fen_from_position(position)
#         else:
#             return False, get_fen_from_position(position)
#
#     def get_engine_move(self, fen):
#         position = get_position_from_fen(fen)
#
#         best_score, best_move = minimax(position, -float('inf'), float('inf'), depth=3)
#         print(best_score, best_move)
#
#         # TODO: handle when best_score is -inf or inf
#         if best_move == None:
#             return fen, []
#         else:
#             position.apply_move(best_move)
#             possible_moves = position.get_all_moves()
#             new_fen = get_fen_from_position(position)
#
#             # print(best_move, f"({best_score/100})")
#             # print()
#             # print(possible_moves)
#
#             return new_fen, possible_moves

import json
import os
import random
from core.application.repositories.unit_of_work import IUnitOfWork
from core.application.schemas.game import GameCreateSchema
from core.domain.fen.MoveParser import MoveParser
from core.domain.fen.fen_generator import get_position_from_fen, get_fen_from_position
from core.domain.minimax.minimax import minimax
from core.domain.engine.Position import Position

class GameService:
    def __init__(self, uow: IUnitOfWork):
        self.uow = uow
        self.openings = self.load_openings()

    def load_openings(self):
        current_dir = os.path.dirname(__file__)
        file_path = os.path.join(current_dir, 'all_openings.json')

        print(f"Loading openings from {file_path}")
        with open(file_path, 'r') as file:
            data = json.load(file)
        return data['openings']

    def create_game(self, game: GameCreateSchema):
        with self.uow as uow:
            game_dict = {
                "created_on": game.created_on,
                "user_id": game.user_id
            }
            game_db = uow.game.add_one(game_dict)
            uow.commit()
            return game_db

    def is_user_move_legal(self, fen, from_square, to_square):
        position = get_position_from_fen(fen)

        piece, color = position.get_piece_and_color_by_square(int(from_square))

        move = {'piece': piece, 'color': color, 'square': int(from_square), 'move': int(to_square)}
        if position.is_move_legal(move):
            position.apply_move(move)
            return True, get_fen_from_position(position)
        else:
            return False, get_fen_from_position(position)

    def get_opening_move(self, fen, selected_variations):
        best_moves = []
        check_for_openings = True
        if len(selected_variations) == 0:
            check_for_openings = False
        for opening in self.openings:
            # somehow filter opening["variations"] ?!
            for variation in opening["variations"]:
                if check_for_openings and variation["name"] in selected_variations:
                    for move in variation['moves']:
                        if move['fen'].split()[:3] == fen.split()[:3]:
                            best_moves.append(move)
                else:
                    for move in variation['moves']:
                        # for variation
                        if move['fen'].split()[:3] == fen.split()[:3]:
                            best_moves.append(move)

        if len(best_moves):
            rand_index = random.randint(0, len(best_moves)-1)
            return best_moves[rand_index]
        return None

    def get_move(self, fen, selected_variations):
        move = self.get_opening_move(fen, selected_variations)
        print(move)
        if move is None:
            fen, possible_moves = self.get_engine_move(fen)
        else:
            move = move['move']
            move_parser = MoveParser()
            move_detail = move_parser.convert_san_to_move_detail(fen, move)
            fen, possible_moves = self.apply_best_move(fen, move_detail)
        return fen, possible_moves

    def apply_best_move(self, fen, move):
        position = get_position_from_fen(fen)
        position.apply_move(move)

        new_fen = get_fen_from_position(position)
        possible_moves = position.get_all_moves()

        return new_fen, possible_moves

    def get_engine_move(self, fen):
        position = get_position_from_fen(fen)
        best_score, best_move = minimax(position, -float('inf'), float('inf'), depth=3)
        if best_move is None:
            return fen, []
        else:
            position.apply_move(best_move)
            possible_moves = position.get_all_moves()
            new_fen = get_fen_from_position(position)
            return new_fen, possible_moves

    def get_all_openings(self):
        # TODO: remove dubliaz code
        current_dir = os.path.dirname(__file__)
        file_path = os.path.join(current_dir, 'all_openings.json')

        print(f"Loading openings from {file_path}")
        with open(file_path, 'r') as file:
            openings = json.load(file)

        transformed_openings = []

        for opening in openings["openings"]:
            print(opening)
            transformed_opening = {
                "name": opening["name"],
                "variations": [{"name": variation["name"]} for variation in opening.get("variations", [])]
            }
            transformed_openings.append(transformed_opening)

        return {"openings": transformed_openings}
