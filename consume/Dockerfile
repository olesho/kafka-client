FROM golang:latest

RUN go get -u github.com/segmentio/kafka-go

WORKDIR /go/src/github.com/olesho/kfk2

ADD . .

WORKDIR /go/src/github.com/olesho/kfk2/consume/

CMD go run main.go
