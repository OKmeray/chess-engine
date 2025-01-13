from core.application.schemas.token import TokenSchema
from core.application.security.auth_handler import encode_jwt
from fastapi import APIRouter, Depends, HTTPException
from core.application.schemas.user import UserCreateSchema, UserLoginSchema
from core.application.services.user_service import UserService
from core.application.repositories.unit_of_work import UnitOfWork, IUnitOfWork
from v1.dependencies import get_uow

router = APIRouter()


def get_user_service(uow: IUnitOfWork = Depends(get_uow)):
    return UserService(uow)


@router.post("/register", response_model=TokenSchema)
def register(user: UserCreateSchema, user_service: UserService = Depends(get_user_service)):
    new_user = user_service.register_user(user)
    if not new_user:
        raise HTTPException(status_code=400, detail="Username already registered")
    access_token = encode_jwt(new_user.id)
    return {"access_token": access_token, "token_type": "bearer"}


@router.post("/login", response_model=TokenSchema)
def login(form_data: UserLoginSchema, user_service: UserService = Depends(get_user_service)):
    access_token = user_service.authenticate_user(form_data.email, form_data.password)
    if not access_token:
        raise HTTPException(status_code=401, detail="Incorrect username or password")
    return {"access_token": access_token, "token_type": "bearer"}


@router.post("/verify-token", response_model=bool)
def verify_token(token: TokenSchema, user_service: UserService = Depends(get_user_service)):
    access_token = token.access_token
    result = user_service.verify_token(access_token)
    return True


@router.post("/verify-token")
def verify_token(token: TokenSchema, user_service: UserService = Depends(get_user_service)):
    print(token)
    is_valid = user_service.verify_token(token.access_token)
    if not is_valid:
        raise HTTPException(status_code=403, detail="Invalid token or expired token.")
    return {"message": "Token is valid."}
