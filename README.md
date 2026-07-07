# Finstat

Веб-приложение для учёта личных финансов

## Состав

Проект состоит из:
1. frontend - Vite React.js frontend.
2. backend - Go backend.
3. migrator - мигратор БД на Go.
4. postgres - SQL БД для хранения данных пользователей.
5. nginx - веб-сервер.

## Сборка

Проект собирается через Docker с помощью комманд:
```
docker compose build
docker compose up
```

После этого в браузере достаточно написать `localhost:8080`
