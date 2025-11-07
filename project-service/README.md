## Project Service

Сервис управления проектами на Go (Gin + GORM + SQLite).

### Возможности
- Создание и управление проектами (только для админов)
- Просмотр собственных проектов (для авторизованных пользователей)
- Обновление статуса и прогресса проекта
- Статистика по проектам пользователя
- JWT-аутентификация и проверка роли admin

### Технологии
- Gin — HTTP API
- GORM + SQLite — база данных
- golang-jwt/jwt — валидация JWT

---

### Быстрый старт
1) Требуется Go 1.21+
2) Переменные окружения:
   - `PORT` (по умолчанию `8082`) — порт HTTP сервера
   - `DB_PATH` (по умолчанию `./project.db`) — путь к файлу SQLite
   - `JWT_SECRET` — секрет для проверки JWT (обязательно)

3) Запуск (dev):
```bash
go run ./main.go
```

4) Сборка и запуск бинаря:
```bash
go build -o project-service
./project-service
```

При старте: создаётся файл БД (если отсутствует) и выполняется автомиграция моделей `Project` и `Log`.

---

### Docker

#### Сборка образа
```bash
docker build -t project-service .
```

#### Запуск контейнера
```bash
docker run -d \
  --name project-service \
  -p 8082:8082 \
  -e PORT=8082 \
  -e DB_PATH=/app/data/project.db \
  -e JWT_SECRET=your-secret-key-here \
  -e AUTH_SERVICE_URL=http://auth-service:8080 \
  -v $(pwd)/data:/app/data \
  project-service
```

#### Использование docker-compose
```bash
# Запуск
docker-compose up -d

# Просмотр логов
docker-compose logs -f project-service

# Остановка
docker-compose down
```

**Важно:** Убедитесь, что `auth-service` доступен по адресу, указанному в `AUTH_SERVICE_URL`.

---

### Структура
- `main.go` — точка входа, инициализация конфига/БД, запуск сервера
- `config/` — загрузка конфига из env (`JWT_SECRET`, `DB_PATH`)
- `db/` — инициализация GORM/SQLite и миграции
- `models/` — модели `Project` и `Log`
- `routes/` — регистрация маршрутов
- `handlers/` — обработчики проектов
- `middleware/` — JWT и проверка роли admin
- `logs/` — отправка событий (используется асинхронно при изменениях)

---

### Требования к JWT
Защищённые эндпоинты ожидают JWT со следующими claim'ами:
- `user_id` — ID пользователя (number)
- `is_admin` — признак роли администратора (boolean)

Секрет для валидации берётся из `JWT_SECRET`.

---

### API
Базовый URL: `http://localhost:8082`

#### Админские эндпоинты (требуют `Authorization: Bearer <token>` и `is_admin=true`)
- POST `/api/projects`
  - Request:
    ```json
    {
      "user_id": 1,
      "title": "New Project",
      "description": "Project description"
    }
    ```
  - Response 201: объект `Project`
    ```json
    {
      "id": 1,
      "user_id": 1,
      "title": "New Project",
      "description": "Project description",
      "progress": 0,
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
    ```

- GET `/api/projects`
  - Response 200: список всех проектов
    ```json
    [
      {
        "id": 1,
        "user_id": 1,
        "title": "New Project",
        "description": "...",
        "progress": 50,
        "status": "in_progress",
        "created_at": "...",
        "updated_at": "..."
      }
    ]
    ```

- PATCH `/api/projects/:id`
  - Request (все поля опциональны):
    ```json
    {
      "title": "Updated Title",
      "description": "Updated description"
    }
    ```
  - Response 200: обновлённый `Project`

- PATCH `/api/projects/:id/progress`
  - Request:
    ```json
    {
      "status": "in_progress",
      "progress": 75
    }
    ```
    - `progress` должен быть от 0 до 100
    - `status` — строка (например: "pending", "in_progress", "completed")
  - Response 200: обновлённый `Project`

- DELETE `/api/projects/:id`
  - Response 200:
    ```json
    { "message": "project deleted" }
    ```

#### Пользовательские эндпоинты (требуют `Authorization: Bearer <token>`)
- GET `/api/me/projects`
  - Response 200: список проектов текущего пользователя
    ```json
    [
      {
        "id": 1,
        "user_id": 1,
        "title": "My Project",
        "description": "...",
        "progress": 50,
        "status": "in_progress",
        "created_at": "...",
        "updated_at": "..."
      }
    ]
    ```

- GET `/api/me/projects/summary`
  - Response 200: статистика по проектам пользователя
    ```json
    {
      "total": 10,
      "completed": 3,
      "in_progress": 7
    }
    ```

#### Ошибки
```json
{ "error": "message" }
```

---

### Примеры cURL
Создание проекта (admin):
```bash
curl -X POST http://localhost:8082/api/projects \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"user_id":1,"title":"New Project","description":"Description"}'
```

Обновление прогресса (admin):
```bash
curl -X PATCH http://localhost:8082/api/projects/1/progress \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status":"in_progress","progress":75}'
```

Мои проекты:
```bash
curl http://localhost:8082/api/me/projects \
  -H "Authorization: Bearer $TOKEN"
```

Статистика проектов:
```bash
curl http://localhost:8082/api/me/projects/summary \
  -H "Authorization: Bearer $TOKEN"
```

---

### Примечания
- Все операции логируются в таблицу `Log` асинхронно
- Значения `user_id` и `is_admin` должны присутствовать в JWT
- Прогресс проекта должен быть в диапазоне 0-100
- Для продакшена запускайте за HTTPS-прокси. Секрет `JWT_SECRET` храните вне репозитория.

---

### Лицензия
MIT (или укажите свою)



