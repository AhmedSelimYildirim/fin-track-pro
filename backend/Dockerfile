# Build aşaması - Tam uyumlu sürüm
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Derleme için gerekli paketler
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN go build -o main ./cmd/api/main.go

# Run aşaması
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY ../.env .
EXPOSE 3000
CMD ["./main"]