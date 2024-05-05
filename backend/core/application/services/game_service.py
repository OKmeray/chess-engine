from core.application.repositories.unit_of_work import IUnitOfWork
from core.application.schemas.game import GameCreateSchema
from core.domain.fen.fen_generator import get_position_from_fen
from core.domain.minimax.minimax import minimax


class GameService:
    def __init__(self, uow: IUnitOfWork):
        self.uow = uow

    def create_game(self, game: GameCreateSchema):
        with self.uow as uow:
            game_dict = {
                "created_on": game.created_on,
                "user_id": game.user_id
            }
            game_db = uow.game.add_one(game_dict)
            uow.commit()
            return game_db

    def get_move_from_position(self, fen):
        position = get_position_from_fen(fen)
        best_score, best_move = minimax(position, -float('inf'), float('inf'), depth=2, maximizing_player=True)
        return best_move
