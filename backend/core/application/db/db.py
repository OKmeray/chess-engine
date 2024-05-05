from sqlalchemy import create_engine
from sqlalchemy.orm import DeclarativeBase, sessionmaker
from core.settings import database_url


engine = create_engine(database_url, echo=True)
SessionLocal = sessionmaker(engine, expire_on_commit=False)


class Base(DeclarativeBase):
    pass


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


def create_tables():
    Base.metadata.create_all(bind=engine)


def drop_tables():
    Base.metadata.drop_all(bind=engine)


if __name__ == "__main__":
    create_tables()
    # drop_tables()

