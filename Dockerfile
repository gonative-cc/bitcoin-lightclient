FROM golang:1.23.1-alpine

WORKDIR /app
EXPOSE 9797

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

CMD ["./main", "data/regtest.json"]
