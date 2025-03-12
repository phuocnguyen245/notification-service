package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"notification-service.com/m/pkg/config"
	"notification-service.com/m/pkg/database"
	"notification-service.com/m/pkg/email"
	"notification-service.com/m/pkg/inapp"
	"notification-service.com/m/pkg/kafka"
	"notification-service.com/m/pkg/notification"
	"notification-service.com/m/pkg/sms"
)

func main() {
	// 1. Load cấu hình từ biến môi trường.
	cfg := config.LoadConfig()

	// 2. Khởi tạo kết nối MongoDB.
	dbClient, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Lỗi kết nối MongoDB: %v", err)
	}
	// Tạo collection reference (ở đây database tên là "notification_db", collection "notifications")
	notificationCollection := dbClient.Database("notification_db").Collection("notifications")
	log.Println("Kết nối MongoDB thành công.")

	// 3. Khởi tạo SMS client từ cấu hình.
	smsClient := sms.NewSMSClientTwilio(cfg.SMSApiAccountSID, cfg.SMSApiAuthToken, "+16366890610")

	emailClient := email.NewEmailClient("smtp.gmail.com", 587, "", "", "")

	// Đăng ký endpoint SSE.
	http.HandleFunc("/sse", inapp.SSEHandler)

	// 4. Khởi tạo Kafka consumer.
	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatalf("Lỗi tạo Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// 5. Consume message từ Kafka (sử dụng partition 0 cho demo, có thể mở rộng cho nhiều partition).
	topic := "notification_events"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Lỗi tạo partition consumer: %v", err)
	}
	defer partitionConsumer.Close()
	log.Println("Kafka consumer đang lắng nghe topic: notification_events")

	dlqProducer := kafka.NewProducer(cfg)

	// 6. Xử lý message nhận được.
	for msg := range partitionConsumer.Messages() {
		log.Printf("Nhận được message Kafka: %s", string(msg.Value))
		// Xử lý từng message trên một goroutine riêng.
		go handleMessage(msg.Value, smsClient, emailClient, inapp.Manager, notificationCollection, dlqProducer)
	}
}

// handleMessage parse message JSON thành Notification struct và gọi xử lý gửi thông báo.
func handleMessage(msg []byte, smsClient *sms.SMSClientTwilio, emailClient *email.EmailClient, inAppClient *inapp.SSEManager, collection database.Collection, dlqProducer sarama.SyncProducer) {
	var n notification.Notification
	if err := json.Unmarshal(msg, &n); err != nil {
		log.Printf("Lỗi parse message: %v", err)
		return
	}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Status = "pending"

	// Lưu thông báo vào MongoDB.
	if err := database.InsertNotification(collection, n); err != nil {
		log.Printf("Lỗi lưu thông báo: %v", err)
		return
	}

	// Xử lý gửi thông báo với cơ chế retry (maxRetries = 3)
	if err := notification.ProcessNotification(n, smsClient, emailClient, inAppClient, 3, dlqProducer, "asd"); err != nil {
		log.Printf("Notification %s gửi thất bại: %v", n.ID, err)
		// Cập nhật trạng thái failed trong DB (bổ sung update ở đây)
		database.UpdateNotificationStatus(collection, n.ID, "failed")
	} else {
		// Cập nhật trạng thái sent trong DB.
		database.UpdateNotificationStatus(collection, n.ID, "sent")
	}
}
