from abc import ABC
from contextlib import contextmanager
from core.application.db.db import SessionLocal
from core.application.repositories.user_repository import UserRepository
from core.application.repositories.game_repository import GameRepository


class IUnitOfWork(ABC):
    def __call__(self):
        raise NotImplementedError

    def __enter__(self):
        raise NotImplementedError

    def __exit__(self, exc_type, exc_val, exc_tb):
        raise NotImplementedError

    async def commit(self):
        raise NotImplementedError

    async def rollback(self):
        raise NotImplementedError


class UnitOfWork(IUnitOfWork):
    def __init__(self, session_factory=SessionLocal):
        self.session_factory = session_factory
        self.session = None
        self.user = None
        self.item = None

    @contextmanager
    def __call__(self):
        self.session = self.session_factory()
        self.user = UserRepository(self.session)
        self.game = GameRepository(self.session)
        try:
            yield self
            self.session.commit()
        except Exception:
            self.session.rollback()
            raise
        finally:
            self.session.close()

    def __enter__(self):
        self.session = self.session_factory()
        self.user = UserRepository(self.session)
        self.item = GameRepository(self.session)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if exc_type:
            self.session.rollback()
        else:
            self.session.commit()
        self.session.close()

    def commit(self):
        self.session.commit()

    def rollback(self):
        self.session.rollback()
