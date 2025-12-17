# Dockerfile для backend сервиса
# Многоступенчатая сборка для минимизации размера образа

# Stage 1: Builder
FROM golang:1.21-alpine AS builder

# Установка необходимых инструментов
RUN apk add --no-cache git make gcc musl-dev

# Установка templ для компиляции шаблонов
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /build

# Копирование go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Компиляция Templ шаблонов
RUN templ generate

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o /build/bin/server \
  ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

# Установка необходимых runtime зависимостей
RUN apk --no-cache add ca-certificates tzdata

# Создание непривилегированного пользователя
RUN addgroup -S rally && adduser -S rally -G rally

WORKDIR /app

# Копирование бинарника из builder stage
COPY --from=builder /build/bin/server /app/server

# Копирование статических файлов и шаблонов
COPY --from=builder /build/web/static /app/web/static
COPY --from=builder /build/web/templates /app/web/templates

# Изменение владельца файлов
RUN chown -R rally:rally /app

# Переключение на непривилегированного пользователя
USER rally

# Exposing порт
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Запуск приложения
CMD ["/app/server"]

