# Используем официальный образ Golang для этапа сборки
FROM golang:1.18 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта в рабочую директорию контейнера
COPY . .

# Загружаем зависимости
RUN GO111MODULE=on go mod tidy

# Собираем приложение
RUN GO111MODULE=on go build -o main cmd/main.go

# Используем минимальный образ для финального контейнера
FROM gcr.io/distroless/base-debian10

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main /main

# Запускаем приложение
CMD ["/main"]
