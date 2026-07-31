# site-monitor

## Сервис для мониторинга сайтов

***

### Запуск для разработки

Из корня проекта:

```
go run ./cmd/monitor -config configs/config.yaml
```

При таком запуске в информации о сборке будет указана версия `dev`.

Для запуска с заданной версией:

```
go run -ldflags "-X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=v1.1.0-beta" ./cmd/monitor -config configs/config.yaml
```

Сайты для мониторинга хранятся в PostgreSQL и добавляются через API (`POST /api/v1/sites`), а не из YAML-конфига.
### Сборка

Версия приложения встраивается в бинарный файл через linker flag `-X`:

```
go build -ldflags "-X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=v1.1.0-beta" -o site-monitor.exe ./cmd/monitor
```

### Docker Compose

```
docker compose up --build
```

При старте сервис `migrate` применяет SQL-миграции, затем поднимается `site-monitor`.

Остановка:

```
docker compose down
```

### Volumes

Список volumes:

```
docker volume ls
```

Остановка контейнеров и удаление volumes проекта:

```
docker compose down -v
```

Удалить конкретный volume вручную:

```
docker volume rm site-monitor_site-monitor-postgres-data
```

Имя volume может отличаться — смотреть `docker volume ls`.

### Makefile

Для упрощения типовых операций в корне проекта есть `Makefile`.

Список команд:

```
make help
```

Основные цели:

| Команда | Описание |
|---------|----------|
| `make build` | локальная сборка приложения |
| `make docker-build` | сборка образов через Docker Compose |
| `make clean` | очистка артефактов сборки |
| `make run` | локальный запуск приложения |
| `make up` | запуск всех сервисов (миграции применяются автоматически) |
| `make down` | остановка сервисов |
| `make restart` | перезапуск сервисов |
| `make deps` | загрузка Go-зависимостей |
| `make fmt` | форматирование кода |
| `make test` | запуск тестов из `internal/tests` |
| `make migrate-up` | накатить следующую миграцию (один шаг) |
| `make migrate-up-head` | накатить все новые миграции |
| `make migrate-down` | откатить последнюю миграцию (один шаг) |
| `make migrate-down-base` | откатить все миграции до нуля |
| `make migrate-version` | текущая версия БД |
| `make db-reset` | пересоздание БД (удаление volume) |
| `make logs` | просмотр логов контейнеров |
| `make ps` | статус контейнеров |
| `make shell` | shell в контейнере приложения |

Версию приложения можно задать через переменную `VERSION` (по умолчанию `dev`):

```
make build VERSION=v1.1.0-beta
make up VERSION=v1.1.0-beta
```

### Миграции БД

Миграции лежат в каталоге `migrations/` и применяются через [goose](https://github.com/pressly/goose).

#### Установка goose

Goose не лежит в `go.mod` приложения — запускается через `go run` с build-тегами (только postgres-драйвер):

```
make migrate-version
```

Эквивалент вручную:

```
go run -tags='no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose@v3.27.2 version
```

#### Подключение к БД

В `.env` две строки подключения:

- `DATABASE_URL` — для приложения в Docker (`postgres:5432`, внутренняя сеть)
- `MIGRATE_DATABASE_URL` — для goose с хоста (`localhost` + `POSTGRES_PORT`)

Пример `.env`:

```
DB_USER=monitor
DB_PASSWORD=monitor
DB_NAME=site_monitor_db
DB_SSLMODE=disable
POSTGRES_PORT=5433

DATABASE_URL=postgres://monitor:monitor@postgres:5432/site_monitor_db?sslmode=disable
MIGRATE_DATABASE_URL=postgres://monitor:monitor@localhost:5433/site_monitor_db?sslmode=disable
```

Если `MIGRATE_DATABASE_URL` не задан, `Makefile` соберёт его автоматически из `DB_*` и `POSTGRES_PORT`.

`make` автоматически подхватывает `.env` из корня проекта.

#### Команды

| Команда | Описание |
|---------|----------|
| `make migrate-up` | накатить следующую миграцию (один шаг) |
| `make migrate-up-head` | накатить все новые миграции |
| `make migrate-down` | откатить последнюю миграцию (один шаг) |
| `make migrate-down-base` | откатить все миграции до нуля |
| `make migrate-version` | показать текущую версию БД |

Примеры:

```
make migrate-up-head
make migrate-up
make migrate-down
make migrate-down-base
make migrate-version
```

`migrate-version` показывает номер последней применённой миграции.

#### Типовой сценарий

1. Поднять стек (PostgreSQL + миграции + приложение):

```
make up
```

Сервис `migrate` накатывает миграции автоматически до старта `site-monitor`.

2. Проверить версию БД (с хоста):

```
make migrate-version
```

3. Полностью пересоздать БД (удалит данные; миграции накатятся в `db-reset`):

```
make db-reset
```

Команды `make migrate-*` нужны для ручного управления схемой с хоста (локальная разработка без пересоздания контейнеров).
#### Создание новой миграции

```
make migrate-up-head
```

или:

```
go run -tags='no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb' \
  github.com/pressly/goose/v3/cmd/goose@v3.27.2 -dir migrations create add_example_index sql
```

После этого отредактируйте созданный файл в `migrations/`: секции `-- +goose Up` и `-- +goose Down`.

### Генерация Swagger-документации

Из корня проекта:

```
swag init -g cmd/monitor/main.go --parseInternal
```
