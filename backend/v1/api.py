from fastapi import APIRouter
from v1.endpoints import users, games

router = APIRouter()

# Include routers from different endpoint modules
router.include_router(users.router, prefix="/users", tags=["users"])
router.include_router(games.router, prefix="/games", tags=["games"])

# You can also add version-specific dependencies or middleware here
