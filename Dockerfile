FROM golang:1.23.1-alpine

ENV DATA_FILE_PATH=./data/regtest.json

WORKDIR /app
EXPOSE 9797

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

CMD ["./main", "${DATA_FILE_PATH}"]
