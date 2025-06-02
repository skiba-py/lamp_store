# Lamp Store

Интернет-магазин ламп с современным веб-интерфейсом и микросервисной архитектурой (Go + React).

## Структура проекта

```
lamp_store/
├── frontend/           # React-приложение (витрина + админка)
├── products_service/   # Сервис товаров (Go)
├── orders_service/     # Сервис заказов (Go)
├── admin_service/      # Сервис админки (Go)
```

## Основные возможности
- Просмотр и поиск товаров
- Оформление заказов
- Админ-панель для управления товарами и заказами
- Загрузка и хранение изображений товаров

## Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/lamp_store.git
cd lamp_store
```

2. Запустите проект с помощью Docker Compose:
```bash
docker compose up --build
```

- Frontend: http://localhost
- Админка: http://localhost/admin

## Сервисы
- **frontend/** — SPA на React (витрина + админка)
- **products_service/** — REST API для товаров (Go)
- **orders_service/** — REST API для заказов (Go)
- **admin_service/** — REST API для админки (Go)

## Документация
- [Frontend README](frontend/README.md)
- [Products Service README](products_service/README.md)
- [Orders Service README](orders_service/README.md)
- [Admin Service README](admin_service/README.md)

## Лицензия
MIT
