# URL Shortener (Golang)

Сервис для сокращения ссылок, построенный по микросервисной архитектуре. Продемонстрирована работа с gRPC, REST, Kafka, Redis и Prometheus

___

## Архитектура

| Сервис              | Назначение                                                                                                                                       | Технологии                                                   |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------ |
| **Auth Service**    | Регистрация, авторизация и выпуск JWT токенов.                                                                                                   | gRPC, GORM, PostgreSQL, bcrypt, JWT, Prometheus metrics.     |
| **Link Service**    | CRUD‑операции над ссылками, генерация уникальных хэшей и публикация событий о создании ссылок в Kafka.                                           | gRPC, GORM, PostgreSQL, Kafka publisher, Prometheus metrics. |
| **Stat Service**    | Учёт кликов по ссылкам и предоставление статистики. Подписывается на Kafka‑события и сохраняет данные в PostgreSQL.               | gRPC, Kafka subscriber, GORM, Prometheus metrics.            |
| **Gateway Service** | API Gateway поверх gRPC‑сервисов. Реализует JWT‑аутентификацию, кеширование ссылок в Redis, rate‑limiting, прометей‑метрики, OpenAPI‑документацию. | chi, Redis, Kafka publisher, Prometheus, OpenAPI.            |

### Все gRPC‑сервисы имеют:
- цепочки middleware (recovery, logging на zap, метрики, в Auth — дополнительно проверка JWT);
- health‑эндпоинты (/healthz);
- /metrics с Prometheus‑метриками.

### Общие компоненты:
- Kafka (pkg/kafka) — пакет для публикации/чтения событий (link.created, link.visited) с поддержкой ретраев, очередей и метрик.
- Proto — определение протобафов (Auth, Link, Stat). Генерация кода через buf.
- Инструменты
    - tools/echoserver — минимальный HTTP‑сервер для отладки.
    - tools/loadtest — простая утилита нагрузочного тестирования.

### Стек технологий
- Язык: Go 1.24.1
- Базы данных: PostgreSQL (у каждого сервиса своя бд)
- Брокер сообщений: Kafka
- Кеш: Redis
- Сетевой слой: gRPC, chi (HTTP)
- ORM: GORM
- Логирование: zap
- Метрики: Prometheus
- Документация: OpenAPI Swagger
- CI‑утилиты: Makefile, buf

___

## Запуск:

### 1. Требования

- Docker и Docker Compose
- Go 1.24+
- Make (опционально; все команды можно выполнить вручную)

### 2. Подготовка ключей JWT

```bash
openssl genpkey -algorithm RSA -out jwt_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in jwt_private.pem -out jwt_public.pem
```
Файлы jwt_private.pem и jwt_public.pem должны лежать в корне репозитория.

### 3. Конфигурация

| Сервис      | Основные переменные                                                                                                                                                                                 |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Auth**    | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_PRIVATE_KEY_PATH`, `JWT_PUBLIC_KEY_PATH`, `AUTH_GRPC_PORT`                                                                          |
| **Link**    | `DB_*` как выше, `LINK_GRPC_PORT`, `KAFKA_ADDR`                                                                                                                                                     |
| **Stat**    | `DB_*`, `STAT_GRPC_PORT`, `KAFKA_ADDR`                                                                                                                                                              |
| **Gateway** | `GATEWAY_PORT`, `AUTH_GRPC_ADDR`, `LINK_GRPC_ADDR`, `STAT_GRPC_ADDR`, `JWT_PUBLIC_KEY_PATH`, `GATEWAY_CORS_ORIGINS`, параметры Redis, rate‑limit и Kafka (см. `internal/gateway/configs/config.go`) |

#### Пример .env для Auth:

```bash
cmd/               # main.go каждого сервиса + Dockerfile
internal/
  auth/            # Auth service
  link/            # Link service
  stat/            # Stat service
  gateway/         # API Gateway
pkg/kafka/         # Kafka
proto/             # .proto схемы
deploy/            # Prometheus, Kafka
tools/             # loadtest и echoserver
```

### 4. Запуск всех сервисов через Docker Compose

```bash
# если есть make
make up
# без make
docker compose up --build
```

Поднимаются:
- 4 микросервиса (auth, link, stat, gateway),
- 3 PostgreSQL базы и миграции,
- Kafka, Redis,
- Prometheus,

URL по умолчанию:
- Gateway: http://localhost:8080
- Prometheus: http://localhost:9090

### 5. Тестирование

Юнит‑тесты
```bash
go test ./...
```

Нагрузочное тестирование
```bash
go run ./tools/loadtest -url http://localhost:8080/ -n 1000 -c 50
```

Преимущества
- Микросервисная архитектура - чёткое разделение ответственности (аутентификация, ссылки, статистика, gateway).
- gRPC + REST - эффективное взаимодействие между сервисами и удобный API Gateway для внешних клиентов.
- Асинхронная обработка - Kafka события при посещении/создании ссылок.
- Кеширование и rate‑limiting - Redis для быстрых редиректов и защиты от злоупотреблений.
- Наблюдаемость - Prometheus‑метрики, health‑чеки, логирование на zap.
- Документация - автосгенерированные gRPC, OpenAPI схема.
- Инструментарий - Makefile, buf, миграции, простые утилиты для e2e проверки и нагрузочного теста.
