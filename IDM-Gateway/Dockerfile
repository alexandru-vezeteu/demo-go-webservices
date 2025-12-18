FROM python:3.11-slim AS base

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY IDM-Gateway/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY IDM/app/proto/idm.proto ./proto/

RUN mkdir -p ./app/proto && \
    python -m grpc_tools.protoc \
    -I./proto \
    --python_out=./app/proto \
    --grpc_python_out=./app/proto \
    ./proto/idm.proto

COPY IDM-Gateway/app/ ./app/

RUN useradd -m -u 1000 gateway && chown -R gateway:gateway /app
USER gateway

EXPOSE 8000

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
