# Repo-Stat: Distributed GitHub Repository Analyzer
Микросервисная система для получения информации о репозиториях GitHub. Проект реализован на Go с соблюдением принципов Чистой архитектуры (Clean Architecture) и использует gRPC для внутрисетевого взаимодействия.

---

## Архитектура

Система разделена на 4 независимых сервиса, развернутых в едином Docker-окружении (монорепозиторий). 

 
Описание сервисов: 

1. API Gateway — точка входа. Принимает REST-запросы, парсит URL, вызывает Processor/Subscriber по gRPC, отдает JSON. Интегрирован Swagger UI.

2. Processor — сервис-посредник. Принимает запрос от Gateway и прозрачно пробрасывает его в Collector.

3. Collector — рабочий сервис. Инкапсулирует логику запросов к GitHub REST API, маппит данные в общие protobuf-структуры.

4. Subscriber — фиктивный сервис, предоставленный для проверки механизмов Health Check (Ping/Pong).

Технический стек 
    Язык: Go 1.25
    Взаимодействие: gRPC & Protocol Buffers (protobuf)
    Веб-сервер: Стандартная библиотека net/http (Go 1.22+ ServeMux) без использования сторонних фреймворков (gin, echo и т.д.)
    Документация: Swagger / OpenAPI (через http-swagger)
    Инфраструктура: Docker, Docker Compose
    Паттерн: Clean Architecture (Слои: Domain, UseCase, Controller, Adapter)
     
---

## Структура проекта (Монорепозиторий) 

Проект использует единую зависимость (go.mod в корне repo-stat/). Общий код вынесен в папку platform.

```text
repo-stat/
├── platform/         # Общая инфраструктура (httpserver, grpcserver, logger, env)
├── proto/            # .proto контракты и сгенерированные pb.go файлы
├── api/              # Сервис API Gateway
├── processor/        # Сервис Processor
├── collector/        # Сервис Collector
├── subscriber/       # Сервис Subscriber
├── go.mod            # Общий модуль
└── go.sum
```
 
 ---
 
## Запуск системы 

Для запуска необходим Docker и Docker Compose. 

 
    Соберите и запустите контейнеры: 

```bash
docker compose up --build
```

---

## API Документация (Swagger) 

После запуска системы Swagger UI доступен по адресу:
    http://localhost:28080/swagger/index.html

Доступные эндпоинты:

1. Проверка здоровья (Health Check) 

Возвращает статус внутренних микросервисов. 

    Endpoint: GET /api/ping
    Success (200):
     
```json
{
  "status": "ok",
  "services": [
    { "name": "processor", "status": "up" },
    { "name": "subscriber", "status": "up" }
  ]
}
``` 
    Degraded (503): Если хотя бы один сервис недоступен, статус меняется на "degraded", а статус сервиса на "down".

2. Информация о репозитории 

Принимает полный URL репозитория GitHub и возвращает базовую информацию. 

    Endpoint: GET /api/repositories/info?url=<github_url>
    Пример запроса: 
    GET /api/repositories/info?url=https://github.com/golang/go
    Success (200):

```json
{
  "full_name": "golang/go",
  "description": "The Go programming language",
  "stars": 125000,
  "forks": 18000,
  "created_at": "2009-11-10T23:00:00Z"
}
```
    Errors:
        400 Bad Request — URL не передан или имеет неверный формат.
        404 Not Found — Репозиторий не существует на GitHub.
        500 Internal Server Error — Ошибка на стороне Collector или GitHub API.
         
---     

## Запуск тестов 

Интеграционные тесты находятся в папке tests/ и обращаются к API Gateway как внешний клиент. 

Выполните команду из корня task3: 
```bash
cd tests && go test -v
```
