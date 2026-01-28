# Bank

## ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ

- `APP_HOST` — Нужна: да. Пример: `:2334`. Описание: адрес и порт, на которых слушает приложение.
- `ALLOW_ORIGIN` — Нужна: да. Пример: `http://127.0.0.1:5137`. Описание: допустимый Origin для CORS.

- `POSTGRES_HOST` — Нужна: да. Пример: `127.0.0.1`. Описание: хост PostgreSQL.
- `POSTGRES_PORT` — Нужна: да. Пример: `5432`. Описание: порт PostgreSQL.
- `POSTGRES_USER` — Нужна: да. Пример: `postgres`. Описание: пользователь базы данных.
- `POSTGRES_PASSWORD` — Нужна: да. Пример: `helloTiver2004`. Описание: пароль пользователя базы данных.
- `POSTGRES_DB` — Нужна: да. Пример: `postgres`. Описание: имя базы данных.
- `POSTGRES_SSLMODE` — Нужна: да. Пример: `disable`. Описание: режим SSL подключения (например `disable` или `require`).
- `POSTGRES_TIMEZONE` — Нужна: да. Пример: `Europe/Moscow`. Описание: часовой пояс БД.
- `POSTGRES_MIGRATIONS_DIR` — Нужна: да. Пример: `/app/migrations`. Описание: директория с миграциями, в контейнере /app/migrations

- `KEYCLOAK_REALM` — Нужна: да. Пример: `schkola`. Описание: realm в Keycloak.
- `KEYCLOAK_AUTH_SERVER` — Нужна: да. Пример: `https://auth.kgb.su`. Описание: URL сервера аутентификации Keycloak.

## Хелсчек
```bash
curl -f <ip>:<port>/ping
```

## Запуск для разработки
- Скачать **air**
    ```bash
    go install github.com/air-verse/air@latest
    ```

- Создать **.env**
    ```bash
    cp .env.example .env
    ```
- Настроить параметры приложения в **.env**

- Запустить **air**
    ```bash
    cd backend
    go mod download
    air
    ```
## Продакшен
- Создать **.env**
    ```bash
    cp .env.example .env
    ```
- Настроить параметры приложения в **.env**
- Запустить докер компоуз
    ```bash
    docker compose up -d
    ```