## Portfolio Service

Сервис управления портфолио на Go (Gin + GORM + SQLite).

### Возможности
- CRUD-операции с портфолио: создание, просмотр, обновление, удаление
- Публичный доступ ко всем эндпоинтам
- Логирование всех операций

### Технологии
- Gin — HTTP API
- GORM + SQLite — база данных

---

### Быстрый старт
1) Требуется Go 1.21+
2) Переменные окружения:
   - `PORT` (по умолчанию `8083`) — порт HTTP сервера
   - `DB_PATH` (по умолчанию `./portfolio.db`) — путь к файлу SQLite

3) Запуск (dev):
```bash
go run ./main.go
```

4) Сборка и запуск бинаря:
```bash
go build -o portfolio-service
./portfolio-service
```

При старте: создаётся файл БД (если отсутствует) и выполняется автомиграция моделей `Portfolio` и `Log`.

---

### Структура
- `main.go` — точка входа, инициализация БД, запуск сервера
- `config/` — загрузка конфига из env (`DB_PATH`)
- `db/` — инициализация GORM/SQLite и миграции
- `models/` — модели `Portfolio` и `Log`
- `routes/` — регистрация маршрутов
- `handlers/` — обработчики портфолио

---

### API
Базовый URL: `http://localhost:8083`

Все эндпоинты публичные (без аутентификации). Рекомендуется добавить защиту для продакшена.

#### Эндпоинты
- POST `/api/portfolio`
  - Request:
    ```json
    {
      "title": "My Project",
      "description": "Project description",
      "image_url": "https://example.com/image.jpg"
    }
    ```
    - `title` — обязательное поле
    - `description` — опциональное поле
    - `image_url` — опциональное поле
  - Response 201: объект `Portfolio`
    ```json
    {
      "id": 1,
      "title": "My Project",
      "description": "Project description",
      "image_url": "https://example.com/image.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
    ```
  - Примечание: создаёт запись в логе

- GET `/api/portfolio`
  - Response 200: список всех портфолио
    ```json
    [
      {
        "id": 1,
        "title": "My Project",
        "description": "Project description",
        "image_url": "https://example.com/image.jpg",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
    ```
  - Примечание: сортировка по дате создания (новые первыми)

- PATCH `/api/portfolio/:id`
  - Request (все поля опциональны):
    ```json
    {
      "title": "Updated Title",
      "description": "Updated description",
      "image_url": "https://example.com/new-image.jpg"
    }
    ```
  - Response 200: обновлённый объект `Portfolio`
  - Примечание: создаёт запись в логе

- DELETE `/api/portfolio/:id`
  - Response 204 (без тела)
  - Примечание: создаёт запись в логе

#### Ошибки
```json
{ "error": "message" }
```

---

### Примеры cURL
Создание портфолио:
```bash
curl -X POST http://localhost:8083/api/portfolio \
  -H 'Content-Type: application/json' \
  -d '{"title":"My Project","description":"Description","image_url":"https://example.com/image.jpg"}'
```

Список портфолио:
```bash
curl http://localhost:8083/api/portfolio
```

Обновление портфолио:
```bash
curl -X PATCH http://localhost:8083/api/portfolio/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Updated Title"}'
```

Удаление портфолио:
```bash
curl -X DELETE http://localhost:8083/api/portfolio/1
```

---

### Примечания
- Все операции логируются в таблицу `Log`
- Эндпоинты не защищены аутентификацией — рекомендуется добавить middleware для продакшена
- Для продакшена запускайте за HTTPS-прокси

---

### Лицензия
MIT (или укажите свою)



