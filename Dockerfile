FROM golang:alpine AS builder

# Установка Git и других зависимостей сборки
RUN apk add git make

# Рабочая директория внутри контейнера
WORKDIR /app

# Копирование исходного кода приложения внутрь контейнера
COPY . .

# Сборка приложения с помощью команды make build
RUN make build

# Контейнер для запуска приложения
FROM alpine:latest AS runner

# Установка зависимостей времени выполнения
RUN apk add busybox ca-certificates

# Рабочая директория внутри контейнера
WORKDIR /app

# Копирование бинарного файла из контейнера сборки в контейнер для запуска
COPY --from=builder /app/build/package/labile-serv .

# Запуск приложения
CMD ["./thief"]
