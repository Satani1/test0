FROM golang:1.20

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server

EXPOSE 8080

CMD ["/app/server"]