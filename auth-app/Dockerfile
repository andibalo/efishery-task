FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod tidy
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authApp ./cmd

CMD ["./authApp"]