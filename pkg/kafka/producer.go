package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"notification-service.com/m/pkg/config"
)

func NewProducer(cfg *config.Config) sarama.SyncProducer {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true
	dlqProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, producerConfig)

	if err != nil {
		log.Fatalf("Lỗi tạo DLQ producer: %v", err)
	}
	defer dlqProducer.Close()

	return dlqProducer
}
