package sender

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_Dispatch_Email(t *testing.T) {
	d := &Dispatcher{
		TestEmail:  "test@example.com",
		TestMobile: "13800000000",
	}
	err := d.Dispatch(context.Background(), ChannelEmail, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestDispatcher_Dispatch_SMS(t *testing.T) {
	d := &Dispatcher{
		TestEmail:  "test@example.com",
		TestMobile: "13800000000",
	}
	err := d.Dispatch(context.Background(), ChannelSMS, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestDispatcher_Dispatch_Inbox(t *testing.T) {
	d := &Dispatcher{
		TestEmail:  "test@example.com",
		TestMobile: "13800000000",
	}
	err := d.Dispatch(context.Background(), ChannelInbox, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestDispatcher_Dispatch_Push(t *testing.T) {
	d := &Dispatcher{
		TestEmail:  "test@example.com",
		TestMobile: "13800000000",
	}
	err := d.Dispatch(context.Background(), ChannelPush, "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestSMSSender_Send(t *testing.T) {
	s := NewSMSSender(SMSConfig{
		Endpoint:  "http://localhost:8080/sms",
		AccessKey: "test-key",
		SignName:  "Test",
		Template:  "SMS_001",
	})
	err := s.Send(context.Background(), "13800000000", "Test", "Test Body")
	assert.Error(t, err)
}

func TestEmailSender_Send(t *testing.T) {
	s := NewEmailSender(SMTPConfig{
		Host:     "localhost",
		Port:     25,
		Username: "test",
		Password: "test",
		From:     "test@example.com",
	})
	err := s.Send(context.Background(), "test@example.com", "Test", "Test Body")
	assert.Error(t, err)
}
