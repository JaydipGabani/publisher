FROM golang:1.19-buster AS builder

WORKDIR /app

COPY . .

RUN go build -o /main

ENTRYPOINT ["/main"]