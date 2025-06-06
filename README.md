# Scheduler

Базовый планировщик задач. Задачей в контексте данного планировщика считается HTTP запрос на определенный URL,
который необходимо выполнить в определенное время.
Для управления таймерами используется TimeWheel.

## 🧾 Структура тела запроса (`POST /task`)

| Поле         | Тип данных         | Обязательное | Описание                                                                 |
|--------------|--------------------|--------------|--------------------------------------------------------------------------|
| `execute_at` | `string` (RFC3339) | Да           | Время срабатывания таймера в формате RFC3339, например `2025-06-05T22:47:00+03:00`. |
| `url`        | `string` (URL)     | Да           | URL, по которому будет отправлен HTTP-запрос по таймеру.                |
| `method`     | `string`           | Да           | HTTP-метод запроса: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`.            |
| `payload`    | `object` или `null`| Нет          | JSON-объект с телом запроса (для методов с телом: `POST`, `PUT`, и т.д.)|

## Request body example

```json
{
  "execute_at": "2025-06-05T22:00:00+03:00",
  "method": "POST",
  "url": "http://localhost:8080/task",
  "payload": {
    "execute_at": "2025-06-05T22:00:00+03:00",
    "method": "GET",
    "url": "http://localhost:8080/task",
    "payload": {
      "message": "incorrect method for this call!"
    }
  }
}
```

## .env fields

```
DB_HOST="хост сервера с бд"
DB_PORT="порт сервера с бд"
DB_USER="юзер бд"
DB_PASSWORD="пароль бд"
DB_NAME="имя бд"
APP_PORT="порт, на котором должно запуститься приложение"
```

## Setup guide

1. Создать .env в корне проекта. Заполнить данными, например:
    ```
    DB_HOST=localhost
    DB_PORT=5434
    DB_USER=user
    DB_PASSWORD=pass
    DB_NAME=scheduler
    APP_PORT=8080
    ```
2. Поднять контейнеры
    ```bash
    docker-compose up -d
    ```