# 1. Build Aşaması
FROM golang:1.24-alpine AS builder

# Derleme için gerekli sistem paketleri
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Bağımlılıkları önceden çekmek (Cache avantajı sağlar)
COPY go.mod go.sum ./
RUN go mod download

# Tüm proje dosyalarını kopyala
COPY . .

# Uygulamayı derle (Statik olarak derlemek alpine için daha güvenlidir)
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api/main.go

# 2. Çalıştırma (Run) Aşaması
FROM alpine:latest

# Çalışma dizinini oluştur
WORKDIR /app

# Sadece derlenmiş binary dosyasını builder'dan kopyala
COPY --from=builder /app/main .

# KRİTİK: .env dosyasını BURAYA KOPYALAMIYORUZ.
# Runtime (çalışma zamanı) sırasında dışarıdan vereceğiz.

# Uygulamanın dinleyeceği port
EXPOSE 3000

# Uygulamayı çalıştır
CMD ["./main"]