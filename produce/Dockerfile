FROM golang:latest

RUN go get -u github.com/segmentio/kafka-go
RUN go get -u github.com/olesho/kafka-client

WORKDIR /go/github.com/olesho/kafka-client/produce
COPY . .

RUN go install github.com/olesho/kafka-client/produce
CMD produce