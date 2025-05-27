# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy web files
COPY --from=builder /app/web ./web

# Create matrices-data directory
RUN mkdir -p /app/matrices-data

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./main"] 