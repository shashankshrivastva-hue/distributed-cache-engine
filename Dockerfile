FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o cache-server cmd/cache-server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/cache-server .
EXPOSE 6379
CMD ["./cache-server"]
