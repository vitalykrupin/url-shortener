# API документация URL Shortener

## Общая информация

Базовый URL: `http://localhost:8080`

Все ответы возвращаются в формате JSON, если не указано иное.

## Аутентификация

Большинство endpoints требуют аутентификации через JWT токен. 
Токен передается в заголовке `Authorization: Bearer <token>`.

Если токен отсутствует, сервер автоматически генерирует новый и устанавливает его в cookie `Token`.

## Endpoints

### Создание короткого URL

#### Из тела запроса

```
POST /
Content-Type: text/plain

https://example.com
```

Ответ:
```
201 Created
Content-Type: text/plain

http://localhost:8080/abc123
```

#### Из JSON

```
POST /api/shorten
Content-Type: application/json

{
  "url": "https://example.com"
}
```

Ответ:
```
201 Created
Content-Type: application/json

{
  "result": "http://localhost:8080/abc123"
}
```

Если URL уже существует:
```
409 Conflict
Content-Type: application/json

{
  "result": "http://localhost:8080/abc123"
}
```

### Создание нескольких коротких URL

```
POST /api/shorten/batch
Content-Type: application/json

[
  {
    "correlation_id": "1",
    "original_url": "https://example1.com"
  },
  {
    "correlation_id": "2",
    "original_url": "https://example2.com"
  }
]
```

Ответ:
```
201 Created
Content-Type: application/json

[
  {
    "correlation_id": "1",
    "short_url": "http://localhost:8080/abc123"
  },
  {
    "correlation_id": "2",
    "short_url": "http://localhost:8080/def456"
  }
]
```

### Получение оригинального URL

```
GET /{alias}
```

Ответ:
```
307 Temporary Redirect
Location: https://example.com
```

Если URL удален:
```
410 Gone
```

Если URL не найден:
```
404 Not Found
```

### Получение всех URL пользователя

```
GET /api/user/urls
Authorization: Bearer <token>
```

Ответ:
```
200 OK
Content-Type: application/json

[
  {
    "short_url": "http://localhost:8080/abc123",
    "original_url": "https://example.com"
  }
]
```

Если у пользователя нет URL:
```
204 No Content
```

Если пользователь не авторизован:
```
401 Unauthorized
```

### Удаление URL пользователя

```
DELETE /api/user/urls
Authorization: Bearer <token>
Content-Type: application/json

["abc123", "def456"]
```

Ответ:
```
202 Accepted
```

### Проверка подключения к базе данных

```
GET /ping
```

Ответ:
```
200 OK
```

Если база данных недоступна:
```
500 Internal Server Error
```

## Ошибки

Все ошибки возвращаются в формате JSON:

```json
{
  "error": "Описание ошибки"
}
```

Коды ошибок:
- 400 Bad Request - Некорректный запрос
- 401 Unauthorized - Не авторизован
- 404 Not Found - Ресурс не найден
- 409 Conflict - Конфликт (URL уже существует)
- 410 Gone - Ресурс удален
- 500 Internal Server Error - Внутренняя ошибка сервера