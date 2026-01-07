from pydantic import BaseModel, EmailStr, Field
from typing import Optional

class RegisterRequest(BaseModel):
    email: EmailStr = Field(..., max_length=255)
    password: str = Field(..., min_length=1, max_length=255)
    role: str = Field(..., pattern="^(client|owner)$")

class RegisterResponse(BaseModel):
    success: bool
    message: str
    user_id: str

class LoginRequest(BaseModel):
    email: EmailStr = Field(..., max_length=255)
    password: str = Field(..., min_length=1, max_length=255)

class LoginResponse(BaseModel):
    success: bool
    token: str
    message: str
    user_id: str
    role: str
    email: str

class VerifyTokenRequest(BaseModel):
    token: str = Field(..., min_length=1, max_length=1000)

class VerifyTokenResponse(BaseModel):
    valid: bool
    email: str
    message: str
    user_id: str
    role: str
    issuer: str
    expires_at: int
    expired: bool
    blacklisted: bool

class RevokeTokenRequest(BaseModel):
    token: str = Field(..., min_length=1, max_length=1000)

class RevokeTokenResponse(BaseModel):
    success: bool
    message: str
