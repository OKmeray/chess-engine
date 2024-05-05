import datetime
from core.application.repositories.unit_of_work import IUnitOfWork
from core.application.schemas.user import UserCreateSchema
from core.application.security.auth_handler import get_password_hash, encode_jwt, verify_password


class UserService:
    def __init__(self, uow: IUnitOfWork):
        self.uow = uow

    def register_user(self, user: UserCreateSchema):
        with self.uow as uow:
            existing_user = uow.user.find_by_username(user.username)
            if existing_user:
                return None
            user_dict = {
                "username": user.username,
                "email": user.email,
                "hashed_password": get_password_hash(user.password),
                "created_on": datetime.datetime.now()
            }
            user_db = uow.user.add_one(user_dict)
            uow.commit()
            return user_db

    def authenticate_user(self, username: str, password: str):
        with self.uow as uow:
            user = uow.user.find_by_username(username)
            if user and verify_password(password, user.hashed_password):
                return encode_jwt(user.id)
            return None
