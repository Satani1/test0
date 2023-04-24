FROM golang:1.20 AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o appTest ./cmd/app.go

EXPOSE 8080

CMD ["./appTest"]