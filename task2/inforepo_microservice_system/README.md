# GitHub Repository Microservice System

**Распределённая система для получения информации о GitHub репозиториях.**  
Эволюционировала из CLI-инструмента в микросервисную архитектуру с использованием **gRPC** и **Clean Architecture**.

---

## Архитектура

Система состоит из двух сервисов:

| Сервис           | Порт     | Тип          | Роль                                                                 |
|------------------|----------|--------------|----------------------------------------------------------------------|
| **Collector**    | `50051`  | gRPC Server  | Инкапсулирует логику взаимодействия с GitHub API                     |
| **API Gateway**  | `8080`   | REST Server  | Принимает HTTP запросы, транслирует их в Collector по gRPC           |

### Структура взаимодействия

```text
Client (HTTP)
    │
    ▼
┌─────────────────┐      gRPC       ┌─────────────────┐      HTTP       ┌────────────┐
│  API Gateway    │ ─────────────►  │   Collector     │ ─────────────►  │ GitHub API │
│  (Port 8080)    │ ◄─────────────  │   (Port 50051)  │ ◄─────────────  │            │
└─────────────────┘      gRPC       └─────────────────┘      HTTP       └────────────┘
```

## Требования

- Go 1.26+
- Docker & Docker Compose
- Protoc (для генерации, если нужно менять proto)

## Быстрый старт (Docker)

Выполните команду в корне проекта:

```bash
docker-compose up --build
```

Сервисы запустятся:

    - API Gateway: http://localhost:8080
    - Collector: localhost:50051 (внутренний порт)

## Использование

### Swagger UI

Откройте в браузере:http://localhost:8080/swagger/index.html

### REST API (Пример запроса)

Запрос:

curl http://localhost:8080/repo/golang/go

Ответ:

{    "name": "go",    "description": "The Go programming language",    "stars": 121000,    "forks": 17300,    "created_at": "2015-01-16T23:15:58Z"}

### Обработка ошибок

    Если репозиторий не найден на GitHub, сервис вернет HTTP 404.
    Если GitHub API недоступен, сервис вернет HTTP 500.