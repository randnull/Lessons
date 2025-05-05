FROM golang:1.23-alpine

WORKDIR /app

COPY ./user_service/go.mod ./user_service/go.sum ./

RUN go mod download

COPY ./user_service/cmd ./cmd
COPY ./user_service/internal ./internal

RUN go build -o /user-service ./cmd/user-app/main.go

CMD ["/user-service"]
