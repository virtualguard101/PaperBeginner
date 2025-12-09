# Development Dockerfile with hot reload
FROM golang:1.22-alpine

WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ .

# Expose port
EXPOSE 8080

# Run with hot reload
CMD ["air", "-c", ".air.toml"]

