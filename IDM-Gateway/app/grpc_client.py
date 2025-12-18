import grpc
from contextlib import asynccontextmanager
from app.config import get_settings
from app.proto import idm_pb2_grpc

class IDMGrpcClient:
    def __init__(self):
        self.settings = get_settings()
        self.address = f"{self.settings.IDM_GRPC_HOST}:{self.settings.IDM_GRPC_PORT}"
        self.channel = None
        self.stub = None

    async def connect(self):
        self.channel = grpc.aio.insecure_channel(self.address)
        self.stub = idm_pb2_grpc.IdentityServiceStub(self.channel)

    async def close(self):
        if self.channel:
            await self.channel.close()

    def get_stub(self):
        return self.stub

grpc_client = IDMGrpcClient()
