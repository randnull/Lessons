FROM golang:alpine

WORKDIR /app

COPY .. .

RUN go build -o /order-service ./cmd/order-app/main.go

CMD ["/order-service"]