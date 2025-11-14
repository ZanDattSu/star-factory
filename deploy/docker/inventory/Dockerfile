# ============================
# Stage 1: Build stage
# ============================

FROM golang:1.24.2-alpine AS builder

# Устанавливаем git — нужен для загрузки зависимостей из приватных и публичных репозиториев
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.work и его контрольную сумму — это позволяет Go правильно связать модули
COPY go.work go.work.sum ./

# Копируем go.mod и go.sum всех зависимых модулей — нужно для кэширования go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код всех модулей — включая зависимости, т.к. они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY order ./order
COPY inventory ./inventory
COPY payment ./payment

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем бинарный файл Inventory-сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/main.go


# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.21.3

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Запускаем приложение под non-root пользователем
USER appuser

COPY --from=builder /app/app-inventory .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

EXPOSE 50051

# Устанавливаем команду запуска — запускаем наш бинарь
ENTRYPOINT ["./app-inventory"]
