FROM --platform=linux/amd64 golang:1.23.3-bullseye AS builder

RUN apt-get update && apt-get install -y \
    pkg-config libvips-dev gcc

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /app/cmd/service

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /main .

# Финальный образ также принудительно берем под linux/amd64
FROM --platform=linux/amd64 debian:bullseye-slim

RUN apt-get update && apt-get install -y libvips

WORKDIR /
COPY --from=builder /main /main
EXPOSE 8081

CMD ["/main"]
