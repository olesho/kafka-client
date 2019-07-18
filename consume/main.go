package main

import (
	"os"

	kafka "github.com/olesho/kafka-client"
)

var topic = os.Getenv("KAFKA_TOPIC")
var addr = os.Getenv("KAFKA_ADDR")

func main() {
	if topic == "" {
		topic = "test-topic"
	}
	if addr == "" {
		addr = "172.17.0.4:9092"
	}

	kafka.Reader([]string{addr}, "test-client", topic)
}
