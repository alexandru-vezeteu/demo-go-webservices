from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
from app.config import get_settings
from app.grpc_client import grpc_client
from app.routers import auth_router

@asynccontextmanager
async def lifespan(app: FastAPI):
    await grpc_client.connect()
    yield
    await grpc_client.close()

app = FastAPI(
    title="IDM Gateway API",
    description="HTTP/REST to gRPC gateway for Identity Management Service",
    version="1.0.0",
    lifespan=lifespan
)

settings = get_settings()
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(auth_router.router)

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "idm-gateway"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.IDM_GATEWAY_HOST,
        port=settings.IDM_GATEWAY_PORT,
        reload=False
    )
