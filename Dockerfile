FROM golang:1.23.3-alpine AS builder

RUN apk update && apk add --no-cache pkgconfig vips-dev gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /app/cmd/service
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /main .

FROM alpine:3.17

RUN apk add --no-cache vips
WORKDIR /
COPY --from=builder /main /main
EXPOSE 8081

CMD ["/main"]
