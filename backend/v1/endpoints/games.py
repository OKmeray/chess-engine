import asyncio

from starlette.websockets import WebSocketDisconnect

from v1.dependencies import get_uow
from fastapi import APIRouter, WebSocket, Depends, HTTPException, Query
from core.application.services.game_service import GameService
from core.application.repositories.unit_of_work import IUnitOfWork
from core.application.security.auth_bearer import JWTBearer


router = APIRouter()


def get_game_service(uow: IUnitOfWork = Depends(get_uow)):
    return GameService(uow)


@router.get("/openings")
def get_openings(game_service: GameService = Depends(get_game_service)):
    return game_service.get_all_openings()
    # {
    #     "openings": [
    #         {
    #             "name": "Ruy Lopez",
    #             "variations": [
    #                 {"name": "Marshall Attack"},
    #                 {"name": "Closed, Martinez Variation"}
    #             ]
    #         },
    #         {
    #             "name": "Caro-Kann",
    #             "variations": [
    #                 {"name": "Advance Variation"},
    #                 {"name": "Advance Variation, Botvinnik-Carls Defense"},
    #                 {"name": "Advance Variation, Short Variation"},
    #                 {"name": "Advance Variation, Tal Variation"}
    #             ]
    #         }
    #     ]
    # }