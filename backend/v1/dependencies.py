from typing import Annotated
from fastapi import Depends
from core.application.repositories.unit_of_work import IUnitOfWork, UnitOfWork

UOWDep = Annotated[IUnitOfWork, Depends(UnitOfWork)]

def get_uow():
    with UnitOfWork() as uow:
        yield uow
