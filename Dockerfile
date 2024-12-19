# Используем конкретную архитектуру и версию golang
FROM --platform=linux/amd64 golang:1.23.3-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache pkgconfig vips-dev gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей для кэширования
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем все файлы проекта
COPY . .

# Устанавливаем рабочую директорию для сборки
WORKDIR /app/cmd/service

# Сборка бинарного файла
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /main .

# Финальный образ
FROM --platform=linux/amd64 alpine:latest

# Устанавливаем необходимые зависимости
RUN apk add --no-cache vips

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарник
COPY --from=builder /main /main

# Копируем watermark/logo.png
COPY watermark/logo.png /app/watermark/logo.png

# Указываем порт
EXPOSE 8081

# Запуск приложения
CMD ["/main"]
