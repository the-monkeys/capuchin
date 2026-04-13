FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Build context is backend/ — migration module lives at backend/migration/
COPY migration/go.mod migration/go.sum ./
RUN go mod download

COPY migration/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrate/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/migrate ./

CMD ["./migrate"]
