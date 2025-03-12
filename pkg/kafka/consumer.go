package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// NewConsumer tạo và trả về Kafka consumer dựa trên danh sách brokers.
func NewConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Printf("Lỗi tạo Kafka consumer: %v", err)
		return nil, err
	}
	log.Printf("Consumer create success")
	return consumer, nil
}
