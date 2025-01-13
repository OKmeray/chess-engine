from pydantic import BaseModel, EmailStr, SecretStr
from datetime import datetime


class UserCreateSchema(BaseModel):
    email: EmailStr
    username: str
    password: str


class UserLoginSchema(BaseModel):
    email: str
    password: str


class UserUpdateSchema(BaseModel):
    email: EmailStr | None = None
    username: str | None = None
    language: str | None = None
    password: SecretStr | None = None

    class Config:
        from_attributes = True


class UserSchema(BaseModel):
    id: int
    username: str
    created_on: datetime

