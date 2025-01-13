import asyncio
import uvicorn
from fastapi import FastAPI, Depends, Query, HTTPException, WebSocket
from fastapi.middleware.cors import CORSMiddleware
from starlette.websockets import WebSocketDisconnect
from core.application.services.game_service import GameService
from v1.api import router as v1_router
from core.settings import SERVER_HOST, SERVER_PORT, LOG_LEVEL, DEBUG_RELOAD, allowed_origins
from v1.endpoints.games import get_game_service


app = FastAPI(title="Chess Engine", version="1.0")


@app.websocket("/ws/move")
async def chess_move(
    websocket: WebSocket,
    game_service: GameService = Depends(get_game_service)
):
    await websocket.accept()
    try:
        while True:
            # Receive FEN from the client
            data = await asyncio.wait_for(websocket.receive_json(), timeout=600)

            # Validate the user's move
            is_legal, new_fen = game_service.is_user_move_legal(data["fen"], data["from"], data["to"])

            # Send the updated position after the user's move back to the client
            await websocket.send_json({
                "fen": new_fen,
                "isUserTurn": not is_legal
            })

            # Compute the engine's move
            if is_legal:
                engine_fen, engine_possible_moves = game_service.get_move(new_fen, data["selectedVariations"])
                print(engine_possible_moves)
                # Send the updated position after the engine's move back to the client
                await websocket.send_json({
                    "fen": engine_fen,
                    "possibleMoves": engine_possible_moves,
                    "isUserTurn": is_legal
                })
    except asyncio.TimeoutError:
        print("No data received. Keeping connection alive.")
        await websocket.send_json({"response": "Keepalive"})
    except WebSocketDisconnect:
        print("Client disconnected")

app.include_router(v1_router, prefix="/api/v1")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # settings.allowed_origins,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE"],
    allow_headers=["*"],
)

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=SERVER_HOST,
        port=SERVER_PORT,
        log_level=LOG_LEVEL,
        reload=DEBUG_RELOAD
    )
