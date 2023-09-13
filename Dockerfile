FROM golang:1.20-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o test_telegram_bot ./cmd/bot/main.go

CMD ["./test_telegram_bot"]