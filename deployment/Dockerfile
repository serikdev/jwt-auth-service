# deployment/Dockerfile

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o jwt-service ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/jwt-service .
COPY .env .

EXPOSE 8080

CMD ["./jwt-service"]
