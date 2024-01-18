# Build stage
FROM golang:1.21.4-alpine3.17 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs configs/
COPY --from=builder /migrations migrations/

ENV TZ=Asia/Tashkent

EXPOSE 8080

ENTRYPOINT [ "./main", "app", "--addr=0.0.0.0:8080" ]