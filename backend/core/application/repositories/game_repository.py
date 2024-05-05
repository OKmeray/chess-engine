from sqlalchemy.orm import Session

from core.domain.models.game import Game
from core.application.repositories.sqlalchemy_repository import SQLAlchemyRepository


class GameRepository(SQLAlchemyRepository):
    model = Game

    def __init__(self, session: Session):
        self.db = session
