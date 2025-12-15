FROM golang:1.23-alpine AS deps

WORKDIR /app

COPY app/go.mod app/go.sum ./

RUN go mod download


FROM deps AS build

WORKDIR /app

COPY app/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /eventManager .


FROM alpine:latest AS runtime

WORKDIR /

COPY --from=build /eventManager /eventManager

CMD ["/eventManager"]
