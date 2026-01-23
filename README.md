# Банк мотивашек

## Как запустить?

1. ручками через postrgress создаем базу данных

2. запускаем server

```bash
go run server.go
```

3. тестирование
```bash
# изменение баланса
curl -X POST http://localhost:8080/change-balance/1 \
     -H "Content-Type: application/json" \
     -d '{"value": 500}'

# получение данных
curl -X GET http://localhost:8080/users/1
```
## TODO 

сделать авторизацию через keycloack
