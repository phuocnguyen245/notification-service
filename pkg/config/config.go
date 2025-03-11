package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config định nghĩa cấu hình ứng dụng.
type Config struct {
	MongoURI     string
	KafkaBrokers []string
	SMSApiKey    string
	SMSApiSecret string
}

// LoadConfig đọc file .env và trả về Config.
func LoadConfig() *Config {
	// Load file .env nếu có.
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	cfg := &Config{
		MongoURI:     os.Getenv("MONGO_URI"),
		SMSApiKey:    os.Getenv("SMS_API_KEY"),
		SMSApiSecret: os.Getenv("SMS_API_SECRET"),
	}

	// Đọc Kafka brokers (chuỗi phân cách bởi dấu phẩy)
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers != "" {
		cfg.KafkaBrokers = strings.Split(kafkaBrokers, ",")
	}

	return cfg
}
