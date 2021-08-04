
FROM golang:1.16-alpine

WORKDIR /app

ENV POSTGRES_CONNECTION="postgresql://doadmin:xb1dzpk16nieox53@anita-postgresql-do-user-7928138-0.b.db.ondigitalocean.com:25060/defaultdb?sslmode=require"

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /go-moviestore

EXPOSE 8080
