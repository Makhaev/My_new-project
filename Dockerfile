FROM golang:1.23-alpine AS build

# Устанавливаем нужные пакеты
RUN apk add --no-cache git build-base

WORKDIR /app

# Копируем файлы зависимостей и качаем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# Финальный минимальный образ
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=build /app/server /app/server

USER nonroot:nonroot
EXPOSE 8082
CMD ["/app/server"]