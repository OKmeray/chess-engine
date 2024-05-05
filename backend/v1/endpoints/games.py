from v1.dependencies import get_uow
from fastapi import APIRouter, WebSocket, Depends
from core.application.services.game_service import GameService
from core.application.repositories.unit_of_work import IUnitOfWork


router = APIRouter()


def get_game_service(uow: IUnitOfWork = Depends(get_uow)):
    return GameService(uow)


@router.websocket("/ws/move")
async def chess_move(websocket: WebSocket, game_service: GameService = Depends(get_game_service)):
    await websocket.accept()
    try:
        while True:
            # Receive FEN from the client
            fen = await websocket.receive_text()
            print("Received FEN:", fen)  # Debug print to check what's received

            # Get the best move from the game service
            best_move = game_service.get_move_from_position(fen)

            # Send the best move back to the client
            await websocket.send_json({
                "fen": fen,
                "bestMove": best_move
            })
    except Exception as e:
        print(f"Error: {str(e)}")
        await websocket.close()

