FROM golang:1.23-alpine

WORKDIR /app

COPY ./auth_service/go.mod ./auth_service/go.sum ./

RUN go mod download

COPY ./auth_service/cmd ./cmd
COPY ./auth_service/internal ./internal
COPY ./auth_service/pkg ./pkg

RUN go build -o /auth-service ./cmd/main_service/main.go

CMD ["/auth-service"]
