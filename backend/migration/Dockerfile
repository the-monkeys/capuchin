FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Build context is backend/ — migration module lives at backend/migration/
COPY migration/go.mod migration/go.sum ./
RUN go mod download

COPY migration/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrate/main.go

# Final Stage — minimal image, non-root user for security
FROM scratch

WORKDIR /app

# Copy passwd so the non-root user exists in scratch
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/migrate ./

USER nobody

CMD ["./migrate"]
