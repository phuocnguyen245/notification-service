package sms

import (
	"log"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSClientTwilio sử dụng Twilio để gửi SMS.
type SMSClientTwilio struct {
	Client     *twilio.RestClient
	FromNumber string
}

// NewSMSClientTwilio khởi tạo SMSClientTwilio từ cấu hình.
func NewSMSClientTwilio(accountSID, authToken, fromNumber string) *SMSClientTwilio {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &SMSClientTwilio{
		Client:     client,
		FromNumber: fromNumber,
	}
}

// SendSMS gửi SMS tới số điện thoại với nội dung message.
func (c *SMSClientTwilio) SendSMS(toNumber, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(c.FromNumber)
	params.SetBody(message)

	resp, err := c.Client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Lỗi gửi SMS: %v", err)
		return err
	}
	log.Printf("SMS gửi thành công với SID: %s", *resp.Sid)
	return nil
}
