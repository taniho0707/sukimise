FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]