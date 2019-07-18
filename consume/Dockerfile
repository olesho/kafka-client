FROM golang:latest

RUN go get -u github.com/segmentio/kafka-go
RUN go get -u github.com/olesho/kafka-client

WORKDIR /go/github.com/olesho/kafka-client/consume
COPY . .

RUN go install github.com/olesho/kafka-client/consume
CMD consume