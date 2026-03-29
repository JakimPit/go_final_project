# --- Этап 1: сборка ---
# golang:1.25-alpine — компактный образ с компилятором Go
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum отдельно — Docker кеширует этот слой,
# и при изменении только кода зависимости не будут скачиваться заново
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем статический бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo_app .

# --- Этап 2: финальный образ ---
# alpine — минимальный Linux (~5 МБ), без Go-тулчейна
FROM alpine:latest

WORKDIR /app

# Копируем только то, что нужно для запуска
COPY --from=builder /app/todo_app .
COPY --from=builder /app/web ./web

# Порт по умолчанию
EXPOSE 7540

# Переменные окружения (можно переопределить при запуске контейнера)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./todo_app"]
