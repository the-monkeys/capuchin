FROM golang:1.26-alpine AS deps

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Build Stage
FROM deps AS builder

COPY . .

# CGO_ENABLED=0 ensures a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Development Stage - pinned air version for reproducible dev builds
FROM deps AS dev

RUN go install github.com/air-verse/air@v1.61.7

CMD ["air", "-c", "air.toml"]

# Final Stage - minimal image, non-root user for security
FROM scratch

WORKDIR /app

# Copy passwd so the non-root user exists in scratch
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/server ./

EXPOSE 8080

USER nobody

CMD ["./server"]
