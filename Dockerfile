# Stage 1: Build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Only copy the compiled binary
COPY --from=builder /app/main .

CMD ["./main"]