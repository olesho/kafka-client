package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy"
)

func Reader(kafkaBrokerUrls []string, kafkaTopic, groupId string, partition int) {
	cfg := kafka.ReaderConfig{
		Partition: partition,
		Brokers:   kafkaBrokerUrls,
		GroupID:   groupId,
		Topic:     kafkaTopic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB

		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,

		QueueCapacity: 1,

		// CommitInterval:  time.Second,
	}
	fmt.Printf("ReaderConfig: %+v\n", cfg)

	reader := kafka.NewReader(cfg)

	for {
		read(reader, partition, kafkaTopic)
		//fetch(reader, partition, kafkaTopic)
	}
}

func read(reader *kafka.Reader, partition int, kafkaTopic string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("deadline exceeded")
		} else {
			log.Println("error while receiving message: ", err.Error())
		}
		return
	} else {
		fmt.Println("READ:", m)

		err = reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println(err)
		}
	}
}

func fetch(reader *kafka.Reader, partition int, kafkaTopic string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	m, err := reader.FetchMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("deadline exceeded")
		} else {
			log.Println("error while receiving message: ", err.Error())
		}
		return
	} else {
		fmt.Println("READ:", m)

		err = reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	var kafkaAddr = os.Getenv("KAFKA_ADDR")
	var kafkaTopic = os.Getenv("TOPIC")
	var groupId = os.Getenv("GROUP_ID")

	var partitions []kafka.Partition

	// try dialing until Kafka is up and connection established
	var conn *kafka.Conn
	for {
		var err error
		time.Sleep(time.Second * 1)
		conn, err = kafka.DialContext(context.Background(), "tcp", kafkaAddr)
		if err != nil {
			log.Println(err)
			continue
		}

		// vers, err := conn.ApiVersions()
		// if err != nil {
		// 	log.Println(err)
		// }
		// fmt.Println("API versions supported:")
		// for _, v := range vers {
		// 	fmt.Printf("%+v\n", v)
		// }

		partitions, err = conn.ReadPartitions(kafkaTopic)
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}
	defer conn.Close()

	for _, p := range partitions {
		fmt.Printf("Partition available: %+v\n", p)
		Reader([]string{kafkaAddr}, kafkaTopic, groupId, p.ID)
	}
}
