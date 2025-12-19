# Build stage
FROM golang:1.25.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and MySQL client tools
RUN apk --no-cache add ca-certificates curl

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port (default 8080)
EXPOSE 8080

# Run the application
CMD ["./main"]
