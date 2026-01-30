<h1 align="center" style="border-bottom: none">
<br>Bank
</h1>
<p align="center">
<a href="https://github.com/silaeder-labs/bank/actions/workflows/backend.yml"><img src="https://github.com/silaeder-labs/bank/actions/workflows/backend.yml/badge.svg" alt="Backend build"/></a>
</p>

## Что это?
Сервис для контроля балансов пользователей

## JWT требования
- Пользователь
Нужен просто sub, и чтоб верефицировался через jwk
- Сервис
Нужно в scope пунктик `payment_create`, пример `profile payment_create email`

## CLI
- Управление безлимитными балансами
```bash
./main -grant-unlimited <uuid>
# OR for docker
docker compose run --rm backend /app/main -grant-unlimited  <uuid>
```
- Забрать права на безлимитный баланс
```bash
./main -revoke-unlimited <uuid>
# OR for docker
docker compose run --rm backend /app/main -revoke-unlimited <uuid>
```

## ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ
| Переменная | Обязательная | Пример | Описание |
|---|---:|---|---|
| `APP_HOST` | да | `:2334` | адрес и порт, на которых слушает приложение |
| `ALLOW_ORIGIN` | да | `http://127.0.0.1:5137` | допустимый Origin для CORS |
| `POSTGRES_HOST` | да | `127.0.0.1` | хост PostgreSQL |
| `POSTGRES_PORT` | да | `5432` | порт PostgreSQL |
| `POSTGRES_USER` | да | `postgres` | пользователь базы данных |
| `POSTGRES_PASSWORD` | да | `helloTiver2004` | пароль пользователя базы данных |
| `POSTGRES_DB` | да | `postgres` | имя базы данных |
| `POSTGRES_SSLMODE` | да | `disable` | режим SSL подключения |
| `POSTGRES_TIMEZONE` | да | `Europe/Moscow` | часовой пояс БД |
| `POSTGRES_MIGRATIONS_DIR` | да | `/app/migrations` | директория с миграциями в контейнере |

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