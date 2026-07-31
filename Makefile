.PHONY: build docker-build clean

ifneq (,$(wildcard ./.env))
include .env
export
endif

BINARY := site-monitor$(shell go env GOEXE)
VERSION ?= dev
LDFLAGS := -X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=$(VERSION)
MIGRATIONS_DIR := migrations
GOOSE_VERSION ?= v3.27.2
GOOSE_TAGS ?= no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb
GOOSE ?= go run -tags=$(GOOSE_TAGS) github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
MIGRATE_DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(POSTGRES_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# 1. Команды сборки

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/monitor

docker-build:
	docker compose pull
	docker compose build --build-arg VERSION=$(VERSION)

clean:
	go clean
ifeq ($(OS),Windows_NT)
	cmd /C "if exist site-monitor.exe del /Q site-monitor.exe"
	cmd /C "if exist site-monitor del /Q site-monitor"
else
	rm -f site-monitor site-monitor.exe
endif

# 2. Команды запуска

.PHONY: run up down restart

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/monitor -config configs/config.yaml

up:
	docker compose build --build-arg VERSION=$(VERSION)
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

# 3. Команды разработки

.PHONY: deps fmt test

deps:
	go mod download

fmt:
	go fmt ./...

test:
	go test ./internal/tests/...

# TODO: make lint — запуск линтера (когда будет настроен)
# lint:
# 	golangci-lint run

# 4. Команды для работы с БД

.PHONY: db-reset migrate-up migrate-up-head migrate-down migrate-down-base migrate-version

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(MIGRATE_DATABASE_URL)" up-by-one

migrate-up-head:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(MIGRATE_DATABASE_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(MIGRATE_DATABASE_URL)" down

migrate-down-base:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(MIGRATE_DATABASE_URL)" down-to 0

migrate-version:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(MIGRATE_DATABASE_URL)" version

db-reset:
	docker compose down -v
	docker compose build --build-arg VERSION=$(VERSION)
	docker compose up -d --wait postgres
	docker compose run --rm migrate
	docker compose up -d

# 5. Вспомогательные команды

.PHONY: logs ps shell help

logs:
	docker compose logs -f

ps:
	docker compose ps

shell:
	docker compose exec site-monitor sh

help:
	@echo Build:
	@echo   make build              - build Go app locally
	@echo   make docker-build       - build images via docker compose
	@echo   make clean              - remove build artifacts
	@echo.
	@echo Run:
	@echo   make run                - run app locally
	@echo   make up                 - start all services (applies migrations)
	@echo   make down               - stop all services
	@echo   make restart            - restart services
	@echo.
	@echo Development:
	@echo   make deps               - download Go modules
	@echo   make fmt                - format code
	@echo   make lint               - TODO: linter
	@echo   make test               - run tests
	@echo.
	@echo Database:
	@echo   make migrate-up         - apply next migration from host (one step)
	@echo   make migrate-up-head    - apply all pending migrations from host
	@echo   make migrate-down       - rollback last migration (one step)
	@echo   make migrate-down-base  - rollback all migrations to base
	@echo   make migrate-version    - show current database version
	@echo   make db-reset           - recreate DB volume and stack (migrations on up)
	@echo.
	@echo Helpers:
	@echo   make logs               - follow container logs
	@echo   make ps                 - show container status
	@echo   make shell              - shell into app container
	@echo   make help               - show this help
	@echo.
	@echo Variables:
	@echo   VERSION=v1.0.0          - app version (default: dev)
