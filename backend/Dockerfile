FROM golang:1.25.5-alpine AS deps

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build Stage
FROM deps AS builder

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Development Stage
FROM deps AS dev

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", "air.toml"]

# Final Stage
FROM scratch

# Set working directory to the app root
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server ./

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]

