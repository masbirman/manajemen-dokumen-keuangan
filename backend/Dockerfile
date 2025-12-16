# backend/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Install air for hot reload (v1.49.0 is compatible with Go 1.22) and git
RUN go install github.com/cosmtrek/air@v1.49.0
RUN apk add --no-cache git

# Copy go mod files first for caching
COPY go.mod ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Download all dependencies
RUN go mod tidy

# Expose port
EXPOSE 8000

# Run with hot reload for development
CMD ["air", "-c", ".air.toml"]
