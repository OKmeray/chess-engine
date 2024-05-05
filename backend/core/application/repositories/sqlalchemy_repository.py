from core.application.repositories.repository import AbstractRepository
from sqlalchemy.orm import Session


class SQLAlchemyRepository(AbstractRepository):
    model = None

    def __init__(self, session: Session):
        self.db = session

    def add_one(self, data):
        db_item = self.model(**data)
        self.db.add(db_item)
        self.db.commit()
        self.db.refresh(db_item)
        return db_item

    def find_all(self, skip: int = 0, limit: int = 100):
        return self.db.query(self.model).offset(skip).limit(limit).all()

    def find_by_id(self, id: int):
        return self.db.query(self.model).filter(self.model.id == id).first()

    def update_one(self, id: int, data):
        db_item = self.db.query(self.model).filter(self.model.id == id).first()
        for var, value in vars(data).items():
            if value is not None:
                setattr(db_item, var, value)
        self.db.commit()
        return db_item

    def delete_one(self, id: int):
        db_item = self.db.query(self.model).filter(self.model.id == id).first()
        if db_item:
            self.db.delete(db_item)
            self.db.commit()
            return db_item
