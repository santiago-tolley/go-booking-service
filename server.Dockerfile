from golang:latest

ARG service

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/server/main.go