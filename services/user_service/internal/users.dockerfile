FROM golang:alpine

WORKDIR /app

COPY .. .

RUN go build -o /user-service ./cmd/user-app/main.go

CMD ["/user-service"]