FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o izdu-client

ENTRYPOINT ["/app/izdu-client"]