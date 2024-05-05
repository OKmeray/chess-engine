from sqlalchemy.orm import Session
from core.application.repositories.sqlalchemy_repository import SQLAlchemyRepository
from core.domain.models.user import User


class UserRepository(SQLAlchemyRepository):
    model = User

    def __init__(self, session: Session):
        self.db = session

    def find_by_email(self, email: str):
        return self.db.query(User).filter(User.email == email).first()

    def find_by_username(self, username: str):
        return self.db.query(User).filter(User.username == username).first()
