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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go

# ===========================
# Stage 2: Final Image
# ===========================
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install golang-migrate
ENV MIGRATE_VERSION=v4.18.2
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin/ \
    && chmod +x /usr/local/bin/migrate \
    && migrate -version

# Copy the built application from the builder container
COPY --from=builder /app/main .

# Copy the directory with migrations
ARG DB_MIGRATIONS_PATH
COPY --from=builder /app/${DB_MIGRATIONS_PATH} ${DB_MIGRATIONS_PATH}

# Copy entrypoint.sh into the container
COPY entrypoint.sh .

# Grant execution permissions to the entrypoint.sh script
RUN chmod +x entrypoint.sh

# Set ENTRYPOINT and pass an argument (path to the executable file)
ENTRYPOINT ["./entrypoint.sh"]