package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(topic, kafkaAddr string, numPartitions int) error {
	conn, err := kafka.Dial("tcp", kafkaAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: 1,
	})
}

func DeleteTopic(topic, kafkaAddr string) error {
	conn, err := kafka.Dial("tcp", kafkaAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.DeleteTopics(topic)
}

func ListPartitions(topic, kafkaAddr string) ([]int, error) {
	conn, err := kafka.Dial("tcp", kafkaAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return nil, err
	}

	res := []int{}
	for _, p := range partitions {
		res = append(res, p.ID)
	}
	return res, nil
}

func Reader(kafkaBrokerUrls []string, clientId, topic string) {
	config := kafka.ReaderConfig{
		Brokers:         kafkaBrokerUrls,
		GroupID:         clientId,
		Topic:           topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("error while receiving message: ", err.Error())
			continue
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}

func Test(topic, kafkaAddr string) error {
	conn, err := kafka.DialLeader(context.TODO(), "tcp", kafkaAddr, topic, 1)
	if err != nil {
		return err
	}
	fmt.Println("WTF")
	return conn.Close()
}

func TestPayload(topic, kafkaAddr string) {
	start := time.Now()

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "1",
	}

	config := kafka.WriterConfig{
		BatchSize:    1,
		Brokers:      []string{kafkaAddr},
		Topic:        topic,
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	visitWriter := kafka.NewWriter(config)

	defer visitWriter.Close()

	messages := []kafka.Message{}

	sent := 0
	startTime, _ := time.Parse("2006-01-02T15:04:05", "2018-01-01T00:00:00")
	for i := 0; i < 100; i++ {

		messages = append(messages, kafka.Message{
			Value: []byte("Test-test-test"),
			Time:  time.Now(),
		})
	}

	startTime = startTime.Add(time.Hour * 2)
	err := visitWriter.WriteMessages(context.Background(), messages...)
	if err != nil {
		panic(err)
	} else {
		sent += len(messages)
		fmt.Println("sent", sent)
	}

	fmt.Println("Done in:", time.Now().Sub(start).Seconds())
}
