# Используем официальный образ Golang
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Используем минимальный образ для выполнения
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Копируем бинарник из образа builder
COPY --from=builder /app/main .

# Указываем, что контейнер будет слушать на порту 8080
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
