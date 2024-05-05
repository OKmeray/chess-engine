from abc import ABC, abstractmethod
from sqlalchemy.orm import Session


class AbstractRepository(ABC):
    @abstractmethod
    def add_one(self, data):
        raise NotImplementedError

    @abstractmethod
    def find_all(self, skip: int = 0, limit: int = 100):
        raise NotImplementedError

    @abstractmethod
    def find_by_id(self, id: int):
        raise NotImplementedError

    @abstractmethod
    def update_one(self, id: int, data):
        raise NotImplementedError

    @abstractmethod
    def delete_one(self, id: int):
        raise NotImplementedError
