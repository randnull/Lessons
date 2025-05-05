FROM golang:1.23-alpine

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o /user-service ./cmd/user-app/main.go

CMD ["/user-service"]