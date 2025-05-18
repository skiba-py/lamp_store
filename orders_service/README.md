# Сервис заказов (Orders Service)

Сервис для управления заказами в системе интернет-магазина.

## Функциональность

- Создание заказов
- Получение информации о заказе
- Получение списка заказов пользователя
- Обновление статуса заказа
- Удаление заказа

## Технологии

- Go 1.22
- PostgreSQL
- Docker
- Chi (HTTP роутер)
- Viper (конфигурация)

## Структура проекта

```
orders_service/
├── cmd/
│   └── api/            # Точка входа приложения
├── internal/
│   ├── domain/         # Доменные модели и интерфейсы
│   ├── repository/     # Реализация хранилища данных
│   ├── service/        # Бизнес-логика
│   └── handler/        # HTTP-обработчики
├── pkg/
│   ├── config/         # Конфигурация приложения
│   └── logger/         # Логирование
├── migrations/         # Миграции базы данных
├── Dockerfile
├── docker-compose.yml
├── config.yaml
└── README.md
```

## Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone https://github.com/skiba/lamp_store.git
cd lamp_store/orders_service
```

2. Запустите сервис с помощью Docker Compose:
```bash
docker compose up -d
```

3. Сервис будет доступен по адресу: http://localhost:8001

## API Endpoints

### Заказы

- `POST /api/orders` - Создание нового заказа
- `GET /api/orders/{id}` - Получение информации о заказе
- `GET /api/orders/user/{userID}` - Получение списка заказов пользователя
- `PUT /api/orders/{id}` - Обновление заказа
- `DELETE /api/orders/{id}` - Удаление заказа

## Разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите тесты:
```bash
go test ./...
```

3. Соберите приложение:
```bash
go build -o orders_service ./cmd/api
```
