FROM golang:latest

RUN go get -u github.com/segmentio/kafka-go

ADD . .

CMD go run *.go
