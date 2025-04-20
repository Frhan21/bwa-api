# Gunakan image Go versi 1.23.6 (alpine untuk ukuran yang lebih kecil)
FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Pertama copy file dependensi
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Copy seluruh kode sumber
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Gunakan distroless image untuk production
FROM gcr.io/distroless/base-debian12

# Copy binary dari builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs

WORKDIR /app

EXPOSE 8000

# Jalankan binary
CMD ["/app/main"]