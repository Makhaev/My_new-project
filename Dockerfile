# Используем официальный образ Go
FROM golang:1.23

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN go build -o main .

# Открываем порт (если нужно)
EXPOSE 8080

# Команда запуска
CMD ["./main"]