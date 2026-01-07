from fastapi import APIRouter, HTTPException
from app.models.auth_models import (
    RegisterRequest, RegisterResponse,
    LoginRequest, LoginResponse,
    VerifyTokenRequest, VerifyTokenResponse,
    RevokeTokenRequest, RevokeTokenResponse
)
from app.grpc_client import grpc_client
from app.proto import idm_pb2
from app.utils.error_handler import handle_grpc_error
import grpc

router = APIRouter(prefix="/api/idm/auth", tags=["Authentication"])

@router.post("/register", response_model=RegisterResponse)
async def register(request: RegisterRequest):
    try:
        stub = grpc_client.get_stub()
        grpc_request = idm_pb2.RegisterRequest(
            email=request.email,
            password=request.password,
            role=request.role
        )
        response = await stub.Register(grpc_request)

        if not response.success:
            raise HTTPException(status_code=400, detail=response.message)

        return RegisterResponse(
            success=response.success,
            message=response.message,
            user_id=response.user_id
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)

@router.post("/login", response_model=LoginResponse)
async def login(request: LoginRequest):
    try:
        stub = grpc_client.get_stub()
        grpc_request = idm_pb2.LoginRequest(
            email=request.email,
            password=request.password
        )
        response = await stub.Login(grpc_request)

        return LoginResponse(
            success=response.success,
            token=response.token,
            message=response.message,
            user_id=response.user_id,
            role=response.role,
            email=response.email
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)

@router.post("/verify", response_model=VerifyTokenResponse)
async def verify_token(request: VerifyTokenRequest):
    try:
        stub = grpc_client.get_stub()
        grpc_request = idm_pb2.VerifyTokenRequest(token=request.token)
        response = await stub.VerifyToken(grpc_request)

        return VerifyTokenResponse(
            valid=response.valid,
            email=response.email,
            message=response.message,
            user_id=response.user_id,
            role=response.role,
            issuer=response.issuer,
            expires_at=response.expires_at,
            expired=response.expired,
            blacklisted=response.blacklisted
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)

@router.post("/revoke", response_model=RevokeTokenResponse)
async def revoke_token(request: RevokeTokenRequest):
    try:
        stub = grpc_client.get_stub()
        grpc_request = idm_pb2.RevokeTokenRequest(token=request.token)
        response = await stub.RevokeToken(grpc_request)

        return RevokeTokenResponse(
            success=response.success,
            message=response.message
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)
