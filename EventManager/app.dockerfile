FROM golang:1.24-alpine AS deps

WORKDIR /app

COPY EventManager/app/go.mod EventManager/app/go.sum ./
COPY IDM/app/go.mod IDM/app/go.sum /IDM/app/

RUN go mod download


FROM deps AS build

WORKDIR /app

COPY EventManager/app/ ./
COPY IDM/app/ /IDM/app/

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /eventManager .


FROM alpine:latest AS runtime

WORKDIR /

COPY --from=build /eventManager /eventManager

CMD ["/eventManager"]
