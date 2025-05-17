# Lamp Store Frontend

Фронтенд часть интернет-магазина ламп, разработанная с использованием React, Vite и Chakra UI.

## Технологии

- React 18
- Vite
- Chakra UI
- React Router
- ESLint
- Docker

## Требования

- Node.js 20 или выше
- npm 9 или выше
- Docker (опционально)

## Установка и запуск

### Локальная разработка

1. Установите зависимости:
```bash
npm install
```

2. Запустите сервер разработки:
```bash
npm run dev
```

Приложение будет доступно по адресу http://localhost:5173

### Сборка для продакшена

```bash
npm run build
```

### Запуск в Docker

1. Соберите Docker образ:
```bash
docker build -t lamp-store-frontend .
```

2. Запустите контейнер:
```bash
docker run -p 80:80 lamp-store-frontend
```

Приложение будет доступно по адресу http://localhost

## Структура проекта

```
frontend/
├── src/              # Исходный код
├── public/           # Статические файлы
├── dist/            # Собранное приложение
├── Dockerfile       # Конфигурация Docker
├── nginx.conf       # Конфигурация Nginx
└── package.json     # Зависимости и скрипты
```

## Скрипты

- `npm run dev` - Запуск сервера разработки
- `npm run build` - Сборка для продакшена
- `npm run lint` - Проверка кода линтером
- `npm run preview` - Предварительный просмотр собранного приложения
