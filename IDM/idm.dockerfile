FROM golang:1.24-alpine AS deps

WORKDIR /app

COPY app/go.mod app/go.sum ./

RUN go mod download


FROM deps AS build

WORKDIR /app

COPY app/ ./

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /idmService .


FROM alpine:latest AS runtime

WORKDIR /

COPY --from=build /idmService /idmService

CMD ["/idmService"]