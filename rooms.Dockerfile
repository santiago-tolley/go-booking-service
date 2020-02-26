from golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/rooms/main.go