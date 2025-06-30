# ===========================
# Stage 1: Builder
# ===========================
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build a static binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/api-gateway/main.go

# ===========================
# Stage 2: Final Image
# ===========================
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built application from the builder container
COPY --from=builder /app/main .

# Run
CMD ["./main"]