# 

# Stage 1: Build the application
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build binary yang statis
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Create the final, smaller image
FROM alpine:latest
# Install sertifikat SSL agar Go bisa membuat koneksi HTTPS
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Salin hanya binary yang sudah di-build dari stage sebelumnya
COPY --from=builder /app/main .

# Kita tidak lagi menyalin .env file di sini

# Expose port yang digunakan oleh aplikasi
EXPOSE 8080
# Perintah untuk menjalankan aplikasi saat container dimulai
CMD ["./main"]