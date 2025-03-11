package notification

import (
	"fmt"
	"log"
	"math"
	"time"

	"notification-service.com/m/pkg/email"
	"notification-service.com/m/pkg/inapp"
	"notification-service.com/m/pkg/sms"
)

type Notification struct {
	ID        string                 `json:"notificationId" bson:"_id"`
	UserID    string                 `json:"userId" bson:"userId"`
	Channel   string                 `json:"channel" bson:"channel"` // sms, email, push, inapp
	Message   string                 `json:"message" bson:"message"`
	Status    string                 `json:"status" bson:"status"` // pending, sent, failed
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updatedAt"`
	MetaData  map[string]interface{} `json:"metaData" bson:"metaData"`
}

// SendNotification chọn kênh gửi dựa trên trường Channel của notification.
func SendNotification(n Notification, smsClient *sms.SMSClientTwilio, emailClient *email.EmailClient, inAppClient *inapp.SSEManager) error {
	switch n.Channel {
	case "sms":
		phone, ok := n.MetaData["phoneNumber"].(string)
		if !ok || phone == "" {
			return fmt.Errorf("phoneNumber không hợp lệ trong metaData")
		}
		return smsClient.SendSMS(phone, n.Message)
	case "email":
		recipient, ok := n.MetaData["email"].(string)
		if !ok || recipient == "" {
			return fmt.Errorf("email không hợp lệ trong metaData")
		}
		subject := "Thông báo từ hệ thống"
		return emailClient.SendEmail(recipient, subject, n.Message)
	case "inapp":
		return inAppClient.SendNotification(n.UserID, n.Message)
	default:
		return fmt.Errorf("kênh %s chưa được hỗ trợ", n.Channel)
	}
}

// ProcessNotification với cơ chế retry.
func ProcessNotification(n Notification, smsClient *sms.SMSClientTwilio, emailClient *email.EmailClient, inAppClient *inapp.SSEManager, maxRetries int) error {
	var sendErr error
	for i := range maxRetries {
		sendErr = SendNotification(n, smsClient, emailClient, inAppClient)
		if sendErr == nil {
			log.Printf("Notification %s gửi thành công.", n.ID)
			return nil
		}
		backoffDuration := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Printf("Retry gửi notification %s sau %v vì lỗi: %v", n.ID, backoffDuration, sendErr)
		time.Sleep(backoffDuration)
	}
	return sendErr
}
