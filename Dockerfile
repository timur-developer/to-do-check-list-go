FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod .

COPY . .

RUN go build -o main ./cmd

FROM alpine

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env .env

CMD ["./main"]