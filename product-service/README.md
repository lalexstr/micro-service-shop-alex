## Product Service

Сервис управления товарами на Go (Gin + GORM + SQLite).

### Возможности
- Публичный список продуктов
- Админские CRUD-операции: создать, получить список, обновить, удалить
- JWT-аутентификация и проверка роли admin для защищённых маршрутов

### Технологии
- Gin — HTTP API
- GORM + SQLite — база данных
- golang-jwt/jwt — валидация JWT

---

### Быстрый старт
1) Требуется Go 1.21+
2) Переменные окружения:
   - `PORT` (по умолчанию `8081`) — порт HTTP сервера
   - `DB_PATH` (по умолчанию `./product.db`) — путь к файлу SQLite
   - `JWT_SECRET` — секрет для проверки JWT (обязательно для защищённых эндпоинтов)

3) Запуск (dev):

```bash
go run ./main.go
```

4) Сборка и запуск бинаря:

```bash
go build -o product-service-service
./product-service-service
```

При старте: создаётся файл БД (если отсутствует) и выполняется автомиграция модели `Product`.

---

### Структура
- `main.go` — точка входа, инициализация конфига/БД, запуск сервера
- `config/` — загрузка конфига из env (`JWT_SECRET`, `DB_PATH`)
- `db/` — инициализация GORM/SQLite и миграции
- `models/` — модель `Product`
- `routes/` — регистрация маршрутов
- `handlers/` — обработчики продуктов
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
Базовый URL: `http://localhost:8081`

#### Публичные
- GET `/api/products/public`
  - Query: `page` (int, по умолч. 1), `size` (int, по умолч. 10)
  - Response 200:
    ```json
    { "items": [ {"id":1,"title":"...","price":100.0,"image_url":"..."} ], "page":1, "size":10, "total": 1, "pages": 1 }
    ```

#### Админ (требует `Authorization: Bearer <token>` и `is_admin=true`)
- POST `/api/products`
  - Request:
    ```json
    { "title": "iPhone", "description": "...", "price": 999.99, "image_url": "https://..." }
    ```
  - Response 201: объект `Product`

- GET `/api/products`
  - Query: `page`, `size`
  - Response 200: пагинированный список как у публичного, но для админов

- PATCH `/api/products/:id`
  - Request (все поля опциональны):
    ```json
    { "title": "New title", "price": 123.45 }
    ```
  - Response 200: обновлённый `Product`

- DELETE `/api/products/:id`
  - Response 204 (без тела)

#### Ошибки
```json
{ "error": "message" }
```

---

### Примеры cURL
Публичный список:
```bash
curl 'http://localhost:8081/api/products/public?page=1&size=10'
```

Создание товара (admin):
```bash
curl -X POST http://localhost:8081/api/products \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"title":"iPhone","description":"14 Pro","price":999.99,"image_url":"https://example.com/iphone.jpg"}'
```

Обновление товара (admin):
```bash
curl -X PATCH http://localhost:8081/api/products/1 \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"price":899.99}'
```

Удаление товара (admin):
```bash
curl -X DELETE http://localhost:8081/api/products/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

### Примечания
- Значения `user_id` и `is_admin` должны присутствовать в JWT. Убедитесь, что ваш Auth сервис их добавляет.
- Для продакшена запускайте за HTTPS-прокси. Секрет `JWT_SECRET` храните вне репозитория.

---

### Лицензия
MIT (или укажите свою)





