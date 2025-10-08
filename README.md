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

- Go 1.25.2 или выше
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

### Локальный запуск основного сервиса

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

### Запуск сервиса авторизации

Сервис авторизации теперь находится в отдельном репозитории `../auth-service`:

```
cd ../auth-service
go run cmd/auth/main.go
```

Сервис авторизации запускается на порту 8082 и предоставляет следующие эндпоинты:
- POST /api/auth/register - регистрация нового пользователя
- POST /api/auth/login - вход пользователя
- GET /api/auth/profile - защищенный эндпоинт для проверки токена

## Использование

### Регистрация пользователя

```
curl -X POST http://localhost:8082/api/auth/register -H "Content-Type: application/json" -d '{"login":"user1","password":"password123"}'
```

### Вход пользователя

```
curl -X POST http://localhost:8082/api/auth/login -H "Content-Type: application/json" -d '{"login":"user1","password":"password123"}'
```

### Создание короткого URL

```
curl -X POST http://localhost:8080/ -d "https://example.com" -H "Authorization: Bearer <token>"
```

### Создание короткого URL через API

```
curl -X POST http://localhost:8080/api/shorten -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"url":"https://example.com"}'
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

Также можно запустить сервис авторизации отдельно:

```
cd ../auth-service
docker build -t auth-service .
docker run -p 8082:8082 auth-service
```

## Клиент (cmd/client)

Клиент поддерживает вход и регистрацию, а также отправку URL для сокращения.

- При запуске предлагается выбрать действие: вход или регистрация.
- После успешной авторизации токен сохраняется в `jwt_token.txt`.

Переменные окружения:

- `AUTH_SERVER_URL` — адрес сервера авторизации (по умолчанию `http://localhost:8082`).
- `AUTH_TOKEN` — токен авторизации (если задан, клиент использует его напрямую).

Быстрый старт клиента:

```
go run cmd/client/main.go
```

Пример регистрации через клиент:

1) Выберите «2. Зарегистрироваться»,
2) Введите логин и пароль,
3) Введите длинный URL для сокращения.

## Лицензия

MIT
