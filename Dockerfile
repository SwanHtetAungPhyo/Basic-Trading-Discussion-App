FROM golang:1.24-alpine

WORKDIR /app
COPY .. /app

RUN go mod tidy
RUN go build -o websocket-binance .

EXPOSE 8081
CMD ["./websocket-binance"]