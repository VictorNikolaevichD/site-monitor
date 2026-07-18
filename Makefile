.PHONY: build docker-build clean

BINARY := site-monitor$(shell go env GOEXE)
VERSION ?= dev
LDFLAGS := -X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=$(VERSION)

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
	go run -ldflags "$(LDFLAGS)" ./cmd/monitor -config configs/sites.yaml

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

.PHONY: db-reset

# TODO: make migrate-up — применение миграций (когда появятся миграции)
# migrate-up:
# 	...

# TODO: make migrate-down — откат миграций (когда появятся миграции)
# migrate-down:
# 	...

db-reset:
	docker compose down -v
	docker compose build --build-arg VERSION=$(VERSION)
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
	@echo   make up                 - start all services
	@echo   make down               - stop all services
	@echo   make restart            - restart services
	@echo.
	@echo Development:
	@echo   make deps               - download Go modules
	@echo   make fmt                - format code
	@echo   make lint               - TODO: linter
	@echo   make test               - run tests (internal/tests)
	@echo.
	@echo Database:
	@echo   make migrate-up         - TODO: apply migrations
	@echo   make migrate-down       - TODO: rollback migrations
	@echo   make db-reset           - recreate DB (remove volume)
	@echo.
	@echo Helpers:
	@echo   make logs               - follow container logs
	@echo   make ps                 - show container status
	@echo   make shell              - shell into app container
	@echo   make help               - show this help
	@echo.
	@echo Variables:
	@echo   VERSION=v1.0.0          - app version (default: dev)
