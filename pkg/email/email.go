package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

// EmailClient sử dụng gomail để gửi email.
type EmailClient struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewEmailClient khởi tạo EmailClient.
func NewEmailClient(host string, port int, username, password, from string) *EmailClient {
	return &EmailClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail gửi email đến người nhận với tiêu đề và nội dung message.
func (c *EmailClient) SendEmail(to, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Lỗi gửi Email: %v", err)
		return err
	}
	log.Printf("Email gửi đến %s thành công.", to)
	return nil
}
