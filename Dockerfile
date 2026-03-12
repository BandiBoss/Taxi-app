# Builder
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY backend/ ./backend

WORKDIR /app/backend

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Runner 
FROM alpine:edge

WORKDIR /app

COPY --from=builder /app/backend/myapp .

EXPOSE 8080

ENTRYPOINT ["/app/myapp"]
