# site-monitor

[![CI](https://github.com/ViktorNikolaevichD/site-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/ViktorNikolaevichD/site-monitor/actions/workflows/ci.yml)

Сервис мониторинга доступности сайтов на Go.

Периодически проверяет URL по HTTP, хранит результаты в PostgreSQL, отдаёт API (REST и gRPC), публикует события в Kafka и шлёт уведомления в Telegram. Снаружи удобнее ходить через REST API Gateway, который проксирует вызовы в monitor по gRPC.

**Стек:** Go · PostgreSQL · Kafka · gRPC · Docker Compose · goose · slog

---

## Возможности

- CRUD сайтов для мониторинга
- Периодическая проверка доступности (интервал и таймаут в конфиге)
- Текущий статус и история проверок
- REST API на monitor + Swagger
- gRPC `MonitorService` (protobuf / buf)
- API Gateway: REST → gRPC, прокидывание `X-Request-ID`
- gRPC interceptors: logging, recovery, request id через metadata
- События проверок в Kafka → notification-сервис → Telegram
- Миграции БД (goose), CI: тест + lint

---

## Архитектура

```text
Клиент (curl / UI)
        │  HTTP REST
        ▼
   API Gateway (:8082)
        │  gRPC + metadata (x-request-id)
        ▼
   site-monitor (:8080 REST, :9090 gRPC)
        │                │
        ▼                ▼
   PostgreSQL          Kafka
                           │
                           ▼
                    notification (:8081)
                           │
                           ▼
                       Telegram
```

| Сервис | Роль |
|--------|------|
| `site-monitor` | Ядро: проверки, БД, REST, gRPC |
| `gateway` | Внешний REST-вход, клиент к gRPC monitor |
| `notification` | Consumer Kafka → Telegram |
| `postgres` | Хранение сайтов и результатов |
| `kafka` | Очередь событий проверок |

---

## Быстрый старт

Требования: Docker + Docker Compose, скопированный `.env`.

```bash
cp .env.example .env
# при необходимости поправь порты и NOTIFICATION_TG_*

make up
```

Стек поднимает Postgres, миграции, Kafka, monitor, gateway, notification.

Проверка:

```bash
curl -s http://localhost:8082/health
curl -s http://localhost:8082/api/v1/sites
```

Остановка:

```bash
make down
```

Полный сброс БД (удалит volume):

```bash
make db-reset
```

---

## Порты (по умолчанию из `.env.example`)

| Порт | Сервис |
|------|--------|
| `8082` | Gateway REST |
| `8080` | Monitor REST (+ Swagger) |
| `9090` | Monitor gRPC |
| `8081` | Notification health |
| `5433` | Postgres на хосте |
| `9092` | Kafka на хосте |

Через gateway удобнее для «как снаружи». Прямой REST monitor — для отладки и Swagger: `http://localhost:8080/swagger/`.

---

## API (кратко)

Базовый URL gateway: `http://localhost:8082`.

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/health` | health gateway |
| `GET` | `/api/v1/sites` | список сайтов |
| `GET` | `/api/v1/sites/{id}` | сайт по id |
| `POST` | `/api/v1/sites` | добавить сайт |
| `DELETE` | `/api/v1/sites/{id}` | удалить |
| `GET` | `/api/v1/sites/{id}/status` | последний статус |
| `GET` | `/api/v1/sites/{id}/history` | история (`limit`, `offset`) |

Примеры:

```bash
# создать
curl -s -X POST http://localhost:8082/api/v1/sites \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","name":"example"}'

# список
curl -s http://localhost:8082/api/v1/sites

# статус / история
curl -s http://localhost:8082/api/v1/sites/<id>/status
curl -s 'http://localhost:8082/api/v1/sites/<id>/history?limit=10&offset=0'
```

Тот же контракт сайтов есть на monitor (`:8080`, префикс `/api/v1/...`). Swagger: `http://localhost:8080/swagger/`.

gRPC (локально, нужен [grpcurl](https://github.com/fullstorydev/grpcurl)):

```bash
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 monitor.v1.MonitorService/GetSites
```

Proto: `proto/monitor/v1/`, генерация: `make proto`.

---

## Локальный запуск без Docker (только monitor)

Нужны поднятые Postgres (и при необходимости Kafka) и актуальный `.env` / `configs/config.yaml`.

```bash
make migrate-up-head   # схема БД с хоста
make run               # go run ./cmd/monitor
```

Сборка с версией:

```bash
make build VERSION=v1.0.0
```

Gateway локально:

```bash
go run ./cmd/gateway -config configs/gateway.yaml
```

---

## Makefile

Список всех целей: `make help`.

| Команда | Описание |
|---------|----------|
| `make up` / `down` / `restart` | Docker Compose стек |
| `make run` | локальный monitor |
| `make build` | бинарник monitor |
| `make docker-build` | сборка образов |
| `make test` / `test-cover` | тесты |
| `make lint` / `fmt` | линтер / формат |
| `make proto` / `proto-lint` | codegen и lint protobuf |
| `make migrate-up-head` | все миграции с хоста |
| `make migrate-down` | откат на один шаг |
| `make migrate-version` | версия схемы |
| `make db-reset` | пересоздать БД и стек |
| `make logs` / `ps` / `shell` | логи / статус / shell в app |

Версия образа/бинарника: `VERSION=v1.0.0 make up`.

---

## Конфиг и переменные

- Приложение: `configs/config.yaml`, gateway: `configs/gateway.yaml`
- Окружение Compose: `.env` (образец — `.env.example`)
- `DATABASE_URL` — для сервисов в Docker-сети
- `MIGRATE_DATABASE_URL` — для goose с хоста (`localhost` + `POSTGRES_PORT`)
- Telegram: `NOTIFICATION_TG_BOT_TOKEN`, `NOTIFICATION_TG_CHAT_ID`

---

## Миграции

SQL в `migrations/`, инструмент — [goose](https://github.com/pressly/goose). В Compose сервис `migrate` накатывает схему до старта monitor.

С хоста:

```bash
make migrate-up-head
make migrate-version
```

---

## CI

GitHub Actions: `.github/workflows/ci.yml` — `go test` (coverage) и `golangci-lint` на push в `main`/`dev` и на pull request.

Локально: `make test`, `make lint`. Правила линтера: `.golangci.yaml`.

---

## Структура репозитория (фрагмент)

```text
cmd/monitor          — HTTP + gRPC сервис мониторинга
cmd/gateway          — REST API Gateway
cmd/notification     — consumer уведомлений
internal/            — домен, use case, handlers, grpc, gateway, …
proto/               — контракт MonitorService
gen/                 — сгенерированный gRPC/protobuf код
migrations/          — SQL-миграции
configs/             — YAML-конфиги
docs/                — Swagger
```

---

## Лицензия

[MIT](LICENSE)
