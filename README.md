## Запуск приложения

1. Сбор контейнера: `docker build -t file_upload:latest .`
2. Запуск: `docker run -d -p 8081:8081 --name file_upload file_upload:latest`