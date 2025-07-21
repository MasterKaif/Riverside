package utils

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer
var brokers = []string{os.Getenv("KAFKA_HOST")}

var STREAM_PUBLISH_TOPIC = "studio-stream"

func InitKafka() {

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true

	var err error
	Producer, err = sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	log.Println("Kafka producer initialized (IBM Sarama v1.45.2)")
}

func SendMessage(topic string , key string , value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := Producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	log.Printf("Message sent to partition %d at offset %d", partition, offset)

	return nil
}
