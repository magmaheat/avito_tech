# Используем официальный образ Golang на базе Debian Bullseye
FROM golang:1.22.3-bullseye AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum, и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходники
COPY . .

# Собираем приложение
RUN go build -o avito_tech ./cmd/avito_tech

# Создаем финальный образ
FROM golang:1.22.3-bullseye

WORKDIR /root/

# Копируем исполняемый файл из builder-контейнера
COPY --from=builder /app/avito_tech .

# Устанавливаем нужные переменные окружения
ENV CONFIG_PATH=/root/config/local.yaml

# Копируем конфиги
COPY config/ /root/config/

# Определяем команду по умолчанию
CMD ["sh", "-c", "sleep 5 && ./avito_tech"]