FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o hello

CMD ["/app/hello"]