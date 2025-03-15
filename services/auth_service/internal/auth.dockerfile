FROM golang:alpine

WORKDIR /app

COPY .. .

RUN go build -o /auth-service ./cmd/main_service/main.go

CMD ["/auth-service"]