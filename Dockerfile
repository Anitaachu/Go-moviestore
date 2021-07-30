
FROM golang:1.16-alpine

WORKDIR /app

ENV POSTGRES_CONNECTION = POSTGRES_CONNECTION

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /go-moviestore

EXPOSE 8080
