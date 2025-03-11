package kafka

import (
	"github.com/IBM/sarama"
)

// NewConsumer tạo và trả về Kafka consumer dựa trên danh sách brokers.
func NewConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer(brokers, config)
}
