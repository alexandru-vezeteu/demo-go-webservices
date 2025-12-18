from pydantic_settings import BaseSettings
from functools import lru_cache

class Settings(BaseSettings):
    # Service configuration
    IDM_GATEWAY_HOST: str = "0.0.0.0"
    IDM_GATEWAY_PORT: int = 8000

    # IDM gRPC service configuration
    IDM_GRPC_HOST: str = "idm-service"
    IDM_GRPC_PORT: int = 50051

    # CORS configuration
    CORS_ORIGINS: list[str] = ["*"]

    class Config:
        env_file = ".env"
        case_sensitive = True

@lru_cache()
def get_settings() -> Settings:
    return Settings()
