# Kodland Community Site

Веб-приложение для сообщества Kodland, позволяющее управлять пользователями и таблицами данных.

## Функциональность

### Пользовательский интерфейс (index.html)

- Просмотр списка пользователей
- Поиск пользователей по имени
- Сортировка пользователей по:
  - Дате добавления
  - Имени
  - Количеству лайков
- Просмотр детальной информации о пользователе
- Возможность ставить лайки пользователям

### Административная панель (admin.html)

- Управление пользователями:
  - Добавление новых пользователей
  - Редактирование существующих пользователей
  - Удаление пользователей
- Авторизация администратора
- Просмотр статистики пользователей

### Управление таблицами (admin-tables.html)

- Создание и редактирование таблиц
- Управление колонками и данными
- Настройка видимости таблиц

### Просмотр таблиц (index-tables.html)

- Просмотр списка доступных таблиц
- Поиск таблиц по названию
- Сортировка таблиц по:
  - Дате создания
  - Названию
  - Автору
- Просмотр детальной информации о таблице

## Технологии

### Фронтенд

- HTML5
- CSS3
- JavaScript (ES6+)
- Bootstrap
- Tailwind

### Бэкенд

- Go 1.22.0
- PostgreSQL
- JWT для аутентификации
- REST API архитектура

## Структура проекта

### Фронтенд

```
kds frontend2/
├── index.html          # Главная страница с пользователями
├── index-tables.html   # Страница просмотра таблиц
├── admin.html          # Административная панель
├── admin-tables.html   # Управление таблицами
└── README.md           # Документация проекта
```

### Бэкенд

```
kds backend/
├── cmd/
│   └── server/         # Точка входа приложения
├── internal/
│   ├── config/         # Конфигурация приложения
│   ├── handler/        # HTTP обработчики
│   ├── middleware/     # Промежуточное ПО
│   ├── model/          # Модели данных
│   ├── repository/     # Работа с базой данных
│   └── service/        # Бизнес-логика
├── pkg/                # Общие пакеты
├── migrations/         # Миграции базы данных
└── go.mod              # Зависимости Go
```

## Установка и запуск

### Фронтенд

1. Клонируйте репозиторий
2. Откройте файлы в веб-браузере
3. Для работы с административной панелью необходимо авторизоваться

### Бэкенд

1. Установите Go 1.22.0 или выше
2. Установите PostgreSQL
3. Создайте базу данных и пользователя
4. Настройте переменные окружения:

   ```
   # Database settings
   POSTGRES_USER=your-username
   POSTGRES_PASSWORD=your-password
   POSTGRES_DB=your-db-name
   POSTGRES_HOST=db
   POSTGRES_PORT=5432
   ```

# Application settings

APP_PORT=8080
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

# Security settings

JWT_SECRET=your-secret-key-here
REFRESH_SECRET=your-secret-key-here

ADMIN_USERNAME=your-username
ADMIN_PASSWORD=your-password

````

5. Соберите контейнер:
```bash
docker-compose up --build
````

6. Если нужно удалить контейнер полностью:
   ```bash
   docker-compose down -v
   ```

## API

Проект использует REST API, доступное по адресу: `http://localhost:8080/api`

### Аутентификация (В разработке)

- `POST /api/login` - Авторизация администратора
- `POST /api/refresh` - Обновление токена

### Пользователи

- `GET /api/users` - Получить список пользователей
- `GET /api/users/{id}` - Получить пользователя
- `POST /api/admin/users` - Создать пользователя
- `PUT /api/admin/users/{id}` - Обновить пользователя
- `DELETE /api/admin/users/{id}` - Удалить пользователя
- `POST /api/users/{id}/like` - Поставить лайк
- `DELETE /api/users/{id}/like` - Убрать лайк

### Таблицы

- `POST /api/tables` - Создать новую таблицу
- `GET /api/tables` - Получить список таблиц
- `GET /api/tables/{id}` - Получить конкретную таблицу
- `PUT /api/tables/{id}` - Обновить таблицу
- `DELETE /api/tables/{id}` - Удалить таблицу

### Колонки таблиц

- `POST /api/tables/{id}/columns` - Добавить колонку
- `PUT /api/tables/{id}/columns/{columnId}` - Изменить колонку
- `DELETE /api/tables/{id}/columns/{columnId}` - Удалить колонку
- `PUT /api/tables/{id}/columns/reorder` - Изменить порядок колонок

### Данные таблиц

- `POST /api/tables/{id}/rows` - Добавить строку
- `PUT /api/tables/{id}/rows/{rowId}` - Обновить строку
- `DELETE /api/tables/{id}/rows/{rowId}` - Удалить строку
- `GET /api/tables/{id}/rows` - Получить данные с пагинацией

## Требования

### Фронтенд

- Современный веб-браузер с поддержкой JavaScript
- Доступ к API серверу
- Для административных функций требуется авторизация

### Бэкенд

- Go 1.22.0 или выше
- PostgreSQL 12 или выше
- 2GB RAM минимум
- 1 CPU минимум

## Лицензия

MIT
