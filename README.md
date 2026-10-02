# site-monitor

[![CI](https://github.com/ViktorNikolaevichD/site-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/ViktorNikolaevichD/site-monitor/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](#)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)](#)
[![Kafka](https://img.shields.io/badge/Kafka-231F20?logo=apachekafka&logoColor=white)](#)
[![gRPC](https://img.shields.io/badge/gRPC-244c5a?logo=grpc&logoColor=white)](#)
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)](#)
[![Telegram](https://img.shields.io/badge/Telegram-26A5E4?logo=telegram&logoColor=white)](#)
[![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF?logo=githubactions&logoColor=white)](#)

Сервис мониторинга доступности сайтов на Go.

Периодически проверяет URL по HTTP, хранит результаты в PostgreSQL, отдаёт API (REST и gRPC), публикует события в Kafka и отправляет уведомления в Telegram. REST API Gateway проксирует внешние запросы в monitor по gRPC.

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

Требования: Docker, Docker Compose, файл `.env`.

```bash
cp .env.example .env
make up
```

Стек: Postgres, миграции, Kafka, monitor, gateway, notification.

```bash
curl -s 'http://localhost:8082/health'
curl -s 'http://localhost:8082/api/v1/sites'
```

```bash
make down
```

```bash
make db-reset
```

---

## Порты

| Порт | Сервис |
|------|--------|
| `8082` | Gateway REST |
| `8080` | Monitor REST, Swagger |
| `9090` | Monitor gRPC |
| `8081` | Notification health |
| `5433` | Postgres (хост) |
| `9092` | Kafka (хост) |

Swagger monitor: `http://localhost:8080/swagger/`.

---

## API

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

```bash
curl -s -X POST 'http://localhost:8082/api/v1/sites' \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","name":"example"}'

curl -s 'http://localhost:8082/api/v1/sites'

curl -s 'http://localhost:8082/api/v1/sites/<id>/status'
curl -s 'http://localhost:8082/api/v1/sites/<id>/history?limit=10&offset=0'
```

Тот же REST по сайтам доступен на monitor (`:8080`, `/api/v1/...`).

gRPC ([grpcurl](https://github.com/fullstorydev/grpcurl)):

```bash
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 monitor.v1.MonitorService/GetSites
```

Proto: `proto/monitor/v1/`. Генерация: `make proto`.

---

## Локальный запуск без Docker

Нужны Postgres (и при необходимости Kafka), `.env` и `configs/config.yaml`.

```bash
make migrate-up-head
make run
```

```bash
make build VERSION=v1.0.0
```

```bash
go run ./cmd/gateway -config configs/gateway.yaml
```

---

## Makefile

`make help` — полный список целей.

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

Версия: `VERSION=v1.0.0 make up`.

---

## Конфиг и переменные

- `configs/config.yaml` — monitor
- `configs/gateway.yaml` — gateway
- `.env` — окружение Compose (образец: `.env.example`)
- `DATABASE_URL` — сервисы в Docker-сети
- `MIGRATE_DATABASE_URL` — goose с хоста (`localhost` + `POSTGRES_PORT`)
- `NOTIFICATION_TG_BOT_TOKEN`, `NOTIFICATION_TG_CHAT_ID` — Telegram

---

## Миграции

SQL в `migrations/`, инструмент — [goose](https://github.com/pressly/goose). В Compose сервис `migrate` применяет схему до старта monitor.

```bash
make migrate-up-head
make migrate-version
```

---

## CI

`.github/workflows/ci.yml`: `go test` (coverage) и `golangci-lint` на push в `main`/`dev` и на pull request.

Локально: `make test`, `make lint`. Правила: `.golangci.yaml`.

---

## Структура репозитория

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
