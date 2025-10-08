# Dockerfile для URL Shortener

# Используем официальный образ Go как базовый
FROM golang:1.25.2-alpine AS builder

# Установка переменных окружения
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка бинарного файла
RUN go build -a -installsuffix cgo -o shortener cmd/shortener/main.go

# Используем минимальный образ Alpine для запуска приложения
FROM alpine:latest

# Установка рабочей директории
WORKDIR /root/

# Копирование бинарного файла из builder образа
COPY --from=builder /app/shortener .

# Создание директории для файлового хранилища
RUN mkdir -p /tmp

# Экспонирование порта
EXPOSE 8080

# Команда запуска приложения
CMD ["./shortener"]