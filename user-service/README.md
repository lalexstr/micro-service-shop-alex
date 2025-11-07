## User Service

Сервис управления пользователями на Go (Gin + GORM + SQLite).

### Возможности
- Список пользователей с фильтрацией (только для админов)
- Изменение роли пользователя (admin/user)
- Просмотр активности пользователя
- JWT-аутентификация и проверка роли admin

### Технологии
- Gin — HTTP API
- GORM + SQLite — база данных
- golang-jwt/jwt — валидация JWT
- godotenv — загрузка переменных окружения

---

### Быстрый старт
1) Требуется Go 1.21+
2) Установите зависимости:
```bash
go mod tidy
```

3) Создайте файл `.env` в корне сервиса:
```env
PORT=8085
DB_PATH=./user.db
JWT_SECRET=your-secret-key-here
JWT_TTL_MIN=60
```

4) Запуск (dev):
```bash
go run ./main.go
```

5) Сборка и запуск бинаря:
```bash
go build -o user-service
./user-service
```

При старте: создаётся файл БД (если отсутствует) и выполняется автомиграция моделей `User` и `Log`.

---

### Структура
- `main.go` — точка входа, инициализация конфига/БД, запуск сервера
- `config/` — загрузка конфига из `.env` (`JWT_SECRET`, `DB_PATH`, `PORT`, `JWT_TTL_MIN`)
- `db/` — инициализация GORM/SQLite и миграции
- `models/` — модели `User` и `Log`
- `routes/` — регистрация маршрутов
- `handlers/` — обработчики пользователей
- `middleware/` — JWT аутентификация и проверка роли admin

---

### Требования к JWT
Защищённые эндпоинты ожидают JWT со следующими claim'ами:
- `user_id` — ID пользователя (number)

Для админских операций требуется, чтобы пользователь имел роль `admin` в базе данных.

Секрет для валидации берётся из `JWT_SECRET`.

---

### API
Базовый URL: `http://localhost:8085`

Все эндпоинты требуют `Authorization: Bearer <token>` и роль `admin`.

#### Админские эндпоинты
- GET `/api/users`
  - Query параметры (все опциональны):
    - `role` — фильтр по роли (admin/user)
    - `email` — поиск по email (LIKE)
    - `from` — фильтр по дате создания (>=)
    - `to` — фильтр по дате создания (<=)
  - Response 200:
    ```json
    [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "user",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
    ```

- PATCH `/api/users/:id/role`
  - Request:
    ```json
    { "role": "admin" }
    ```
    Допустимые значения: `"admin"` или `"user"`
  - Response 200:
    ```json
    { "id": 1, "role": "admin" }
    ```
  - Примечание: создаёт запись в логе активности

- GET `/api/users/:id/activity`
  - Query параметры (опциональны):
    - `from` — фильтр по дате (>=)
    - `to` — фильтр по дате (<=)
  - Response 200:
    ```json
    [
      {
        "id": 1,
        "user_id": 1,
        "action": "changed role for user id=2 to admin",
        "timestamp": "2024-01-01T00:00:00Z"
      }
    ]
    ```

#### Ошибки
```json
{ "error": "message" }
```

---

### Примеры cURL
Список пользователей:
```bash
curl 'http://localhost:8085/api/users?role=user&email=john' \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

Изменение роли:
```bash
curl -X PATCH http://localhost:8085/api/users/2/role \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"role":"admin"}'
```

Активность пользователя:
```bash
curl 'http://localhost:8085/api/users/1/activity?from=2024-01-01' \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

### Примечания
- Все операции логируются в таблицу `Log`
- Для работы сервиса требуется, чтобы JWT содержал `user_id`
- Роль пользователя проверяется в базе данных при каждом запросе
- Для продакшена запускайте за HTTPS-прокси. Секрет `JWT_SECRET` храните вне репозитория.

---

### Лицензия
MIT (или укажите свою)



