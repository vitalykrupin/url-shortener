# URL Shortener

URL Shortener - это сервис для сокращения URL, написанный на Go.

## Описание

Этот сервис позволяет пользователям создавать короткие ссылки для длинных URL. 
Пользователи могут регистрироваться, создавать короткие URL, просматривать свои URL 
и удалять их при необходимости.

## Особенности

- Создание коротких URL
- Получение оригинального URL по короткому alias
- Просмотр всех URL пользователя
- Удаление URL пользователя
- Поддержка базы данных PostgreSQL
- Поддержка файлового хранилища
- Авторизация через JWT
- Сжатие ответов через gzip
- Логирование запросов

## Требования

- Go 1.19 или выше
- PostgreSQL (опционально)
- Docker (опционально)

## Установка

1. Клонируйте репозиторий:
   ```
   git clone https://github.com/vitalykrupin/url-shortener.git
   ```

2. Перейдите в директорию проекта:
   ```
   cd url-shortener
   ```

3. Установите зависимости:
   ```
   go mod tidy
   ```

## Конфигурация

Сервис можно сконфигурировать с помощью флагов командной строки или переменных окружения:

| Флаг | Переменная окружения | Описание | Значение по умолчанию |
|------|---------------------|----------|----------------------|
| -a | SERVER_ADDRESS | Адрес сервера | localhost:8080 |
| -b | BASE_URL | Базовый URL для ответов | http://localhost:8080 |
| -f | FILE_STORAGE_PATH | Путь к файлу хранилища | /tmp/short-url-db.json |
| -d | DATABASE_DSN | Строка подключения к базе данных | "" |

## Запуск

### Локальный запуск

```
go run cmd/shortener/main.go
```

### Запуск с параметрами

```
go run cmd/shortener/main.go -a localhost:8080 -b http://localhost:8080
```

### Запуск с использованием базы данных

```
go run cmd/shortener/main.go -d "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

## Использование

### Создание короткого URL

```
curl -X POST http://localhost:8080/ -d "https://example.com"
```

### Создание короткого URL через API

```
curl -X POST http://localhost:8080/api/shorten -H "Content-Type: application/json" -d '{"url":"https://example.com"}'
```

### Получение оригинального URL

```
curl -X GET http://localhost:8080/abc123
```

### Получение всех URL пользователя

```
curl -X GET http://localhost:8080/api/user/urls -H "Authorization: Bearer <token>"
```

### Удаление URL пользователя

```
curl -X DELETE http://localhost:8080/api/user/urls -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '["abc123", "def456"]'
```

## Тестирование

Для запуска тестов выполните:

```
go test ./...
```

## Docker

Для запуска в Docker:

```
docker-compose up
```

## Лицензия

MIT
