from pydantic import BaseModel, EmailStr, Field
from typing import Optional

class LoginRequest(BaseModel):
    email: EmailStr
    password: str = Field(..., min_length=1)

class LoginResponse(BaseModel):
    success: bool
    token: str
    message: str
    user_id: str
    role: str
    email: str

class VerifyTokenRequest(BaseModel):
    token: str = Field(..., min_length=1)

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
    token: str = Field(..., min_length=1)

class RevokeTokenResponse(BaseModel):
    success: bool
    message: str
