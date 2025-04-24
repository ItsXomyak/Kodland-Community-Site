# Этап сборки
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Финальный этап
FROM alpine:3.19

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache postgresql-client && \
    rm -rf /var/cache/apk/*

# Копируем бинарный файл и SQL скрипты
COPY --from=builder /app/main .
COPY --from=builder /app/sql/init.sql ./sql/

# Устанавливаем переменные окружения
ENV PGPASSWORD=password

# Команда для запуска
CMD ["sh", "-c", "until pg_isready -h db -p 5432 -U postgres -d kodland; do sleep 1; done && psql -h db -p 5432 -U postgres -d kodland -f /app/sql/init.sql && /app/main"]