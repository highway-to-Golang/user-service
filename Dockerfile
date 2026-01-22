# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make protobuf protobuf-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate proto files
RUN make proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/user-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/user-service .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080 8081 9090

CMD ["./user-service"]
