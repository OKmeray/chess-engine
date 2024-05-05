from pydantic import BaseModel
from datetime import datetime


class GameCreateSchema(BaseModel):
    created_on: datetime
    user_id: int


class GameUpdateSchema(BaseModel):
    moves: str

    class Config:
        from_attributes = True


class ItemSchema(BaseModel):
    id: int
    moves: str
    created_on: datetime
    user_id: int

    class Config:
        from_attributes = True

