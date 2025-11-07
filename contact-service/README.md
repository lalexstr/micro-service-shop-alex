## Contact Service

Сервис управления контактными запросами на Go (Gin + GORM + SQLite).

### Возможности
- Публичная отправка контактных запросов
- Админское управление запросами: просмотр, изменение статуса, удаление
- Фильтрация запросов по статусу
- Логирование всех операций

### Технологии
- Gin — HTTP API
- GORM + SQLite — база данных
- godotenv — загрузка переменных окружения

---

### Быстрый старт
1) Требуется Go 1.21+
2) Установите зависимости:
```bash
go mod tidy
```

3) Создайте файл `.env` в корне сервиса (опционально):
```env
PORT=8084
DB_PATH=./contact.db
LOG_LEVEL=info
```

4) Запуск (dev):
```bash
go run ./cmd/main.go
```

5) Сборка и запуск бинаря:
```bash
go build -o contact-service ./cmd/main.go
./contact-service
```

При старте: создаётся файл БД (если отсутствует) и выполняется автомиграция моделей `ContactRequest` и `Log`.

---

### Структура
- `cmd/main.go` — точка входа, инициализация конфига/БД, запуск сервера
- `config/` — загрузка конфига из `.env` (`PORT`, `DB_PATH`, `LOG_LEVEL`)
- `db/` — инициализация GORM/SQLite и миграции
- `models/` — модели `ContactRequest` и `Log`
- `routes/` — регистрация маршрутов
- `handlers/` — обработчики контактных запросов и логов

---

### Статусы запросов
- `new` — новый запрос (по умолчанию)
- `answered` — запрос обработан
- `deleted` — запрос удалён

---

### API
Базовый URL: `http://localhost:8084`

#### Публичные эндпоинты
- POST `/api/contact-requests`
  - Request:
    ```json
    {
      "contact": "email@example.com или +1234567890"
    }
    ```
  - Response 200:
    ```json
    {
      "message": "Ваш запрос успешно отправлен! Мы свяжемся с вами в ближайшее время."
    }
    ```
  - Примечание: создаёт запись в логе

#### Админские эндпоинты (без аутентификации, но должны быть защищены в продакшене)
- GET `/api/contact-requests/admin`
  - Query параметры (опциональны):
    - `status` — фильтр по статусу (new/answered/deleted)
  - Response 200:
    ```json
    [
      {
        "id": 1,
        "contact": "email@example.com",
        "status": "new",
        "admin_id": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
    ```

- PATCH `/api/contact-requests/admin/:id/status`
  - Request:
    ```json
    {
      "status": "answered",
      "admin_id": 1
    }
    ```
    - `status` — обязательное поле (new/answered/deleted)
    - `admin_id` — опциональное поле (ID администратора, обработавшего запрос)
  - Response 200: обновлённый объект `ContactRequest`
  - Примечание: создаёт запись в логе

- DELETE `/api/contact-requests/admin/:id`
  - Response 204 (без тела)
  - Примечание: создаёт запись в логе

#### Логи
- GET `/api/logs` (если реализовано)
  - Response 200: список логов

#### Ошибки
```json
{ "error": "message" }
```

---

### Примеры cURL
Создание контактного запроса:
```bash
curl -X POST http://localhost:8084/api/contact-requests \
  -H 'Content-Type: application/json' \
  -d '{"contact":"email@example.com"}'
```

Список запросов (admin):
```bash
curl 'http://localhost:8084/api/contact-requests/admin?status=new'
```

Изменение статуса (admin):
```bash
curl -X PATCH http://localhost:8084/api/contact-requests/admin/1/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"answered","admin_id":1}'
```

Удаление запроса (admin):
```bash
curl -X DELETE http://localhost:8084/api/contact-requests/admin/1
```

---

### Примечания
- Все операции логируются в таблицу `Log`
- Админские эндпоинты не защищены JWT в текущей реализации — рекомендуется добавить middleware для продакшена
- Для продакшена запускайте за HTTPS-прокси

---

### Лицензия
MIT (или укажите свою)



