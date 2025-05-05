FROM golang:1.23-alpine

WORKDIR /app

COPY ./order_service/go.mod ./order_service/go.sum ./

RUN go mod download

COPY ./ban_words.txt ./

COPY ./order_service/cmd ./cmd
COPY ./order_service/internal ./internal

RUN go build -o /order-service ./cmd/order-app/main.go

CMD ["/order-service"]