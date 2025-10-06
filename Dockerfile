FROM golang:1.21-alpine

# Установка зависимостей
RUN apk add --no-cache git

# Создание директории для приложения
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка приложения
RUN go build -o server main.go

# Установка переменной окружения (можно переопределить)
ENV SMS_API_KEY="changeme"

# Открытие порта
EXPOSE 8082

# Запуск приложения
CMD ["./server"]