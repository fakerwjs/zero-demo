package sender

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig SMTP 邮件网关配置。
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// EmailSender 通过 SMTP 真实发送邮件。
type EmailSender struct {
	cfg SMTPConfig
}

func NewEmailSender(cfg SMTPConfig) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (s *EmailSender) Send(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	var msg strings.Builder
	msg.WriteString("From: " + s.cfg.From + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg.String()))
}
