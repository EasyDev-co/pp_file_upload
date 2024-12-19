# Используем конкретную архитектуру и версию golang
FROM --platform=linux/amd64 golang:1.23.3-bullseye AS builder

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y \
    pkg-config libvips-dev gcc

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
FROM --platform=linux/amd64 debian:bullseye-slim

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libvips

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарник
COPY --from=builder /main /main

# Копируем .env файл
COPY .env /app/.env

# Копируем watermark/logo.png
COPY watermark/logo.png /app/watermark/logo.png
COPY /etc/letsencrypt/live/upload.photodetstvo.ru/* /etc/letsencrypt/live/upload.photodetstvo.ru/

# Указываем порт
EXPOSE 8081

# Запуск приложения
CMD ["/main"]
