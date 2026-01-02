FROM golang:1.24-alpine AS deps

WORKDIR /app

COPY User/app/go.mod User/app/go.sum ./
COPY IDM/app/go.mod IDM/app/go.sum /IDM/app/
RUN go mod download


FROM deps AS build

WORKDIR /app

COPY User/app/ ./
COPY IDM/app/ /IDM/app/

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /userService .


FROM alpine:latest AS runtime

WORKDIR /

COPY --from=build /userService /userService

CMD ["/userService"]