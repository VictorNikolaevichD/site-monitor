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

### Docker

Сборка образа из корня проекта:

```
docker build -t site-monitor .
```

Запуск контейнера:

```
docker run --rm -p 8080:8080 -v ./configs/sites.yaml:/config/sites.yaml:ro site-monitor
```

### Генерация Swagger-документации

Из корня проекта:

```
swag init -g cmd/monitor/main.go --parseInternal
```
