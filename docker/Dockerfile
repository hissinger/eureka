FROM golang:1.14-buster

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o client main.go

CMD ["./client"]