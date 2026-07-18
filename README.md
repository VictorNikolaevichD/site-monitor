# site-monitor

## Сервис для мониторинга сайтов

***

### Запуск для разработки

Из корня проекта:

```
go run ./cmd/monitor -config configs/sites.yaml
```

При таком запуске в информации о сборке будет указана версия `dev`.

Для запуска с заданной версией:

```
go run -ldflags "-X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=v1.1.0-beta" ./cmd/monitor -config configs/sites.yaml
```

### Сборка

Версия приложения встраивается в бинарный файл через linker flag `-X`:

```
go build -ldflags "-X gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo.Version=v1.1.0-beta" -o site-monitor.exe ./cmd/monitor
```

### Docker Compose

```
docker compose up --build
```

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
| `make up` | запуск всех сервисов |
| `make down` | остановка сервисов |
| `make restart` | перезапуск сервисов |
| `make deps` | загрузка Go-зависимостей |
| `make fmt` | форматирование кода |
| `make test` | запуск тестов из `internal/tests` |
| `make db-reset` | пересоздание БД (удаление volume) |
| `make logs` | просмотр логов контейнеров |
| `make ps` | статус контейнеров |
| `make shell` | shell в контейнере приложения |

Версию приложения можно задать через переменную `VERSION` (по умолчанию `dev`):

```
make build VERSION=v1.1.0-beta
make up VERSION=v1.1.0-beta
```

### Генерация Swagger-документации

Из корня проекта:

```
swag init -g cmd/monitor/main.go --parseInternal
```
