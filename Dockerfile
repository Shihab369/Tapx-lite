# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download dependencies first (layer caching optimization).
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically linked binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tapx-server ./cmd/server

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS calls.
RUN apk --no-cache add ca-certificates

# Copy only the binary from the builder stage.
COPY --from=builder /app/tapx-server .

EXPOSE 8080

CMD ["./tapx-server"]