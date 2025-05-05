FROM golang:1.23-alpine

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o /order-service ./cmd/order-app/main.go

CMD ["/order-service"]