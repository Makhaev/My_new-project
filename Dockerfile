FROM golang:1.21-alpine

# Установим зависимости для CGO и SQLite
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Включаем CGO
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Сборка
RUN go build -o server main.go

# Открываем порт
EXPOSE 8082

# Запуск
CMD ["./server"]
