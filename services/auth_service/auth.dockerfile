FROM golang:1.23-alpine

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o /auth-service ./cmd/main_service/main.go

CMD ["/auth-service"]