FROM golang:1.23-alpine AS deps

WORKDIR /app

COPY app/go.mod app/go.sum ./
COPY ../IDM/app/go.mod ../IDM/app/go.sum /IDM/app/
RUN go mod download


FROM deps AS build

WORKDIR /app

COPY app/ ./
COPY ../IDM/app/ /IDM/app/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /userService .


FROM alpine:latest AS runtime

WORKDIR /

COPY --from=build /userService /userService

CMD ["/userService"]