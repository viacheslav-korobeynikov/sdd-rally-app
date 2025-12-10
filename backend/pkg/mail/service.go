package mail

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/viacheslav-korobeynikov/sdd-rally-app/config"
)

type SMTPClient struct {
	Host string
	Port string
	From string
	Auth smtp.Auth
}

func NewSMTPClient() (*SMTPClient, error) {
	// Получаем конфиг SMTP
	cfg := config.NewSMTPConfig()
	// Проверяем, что хост и порт не пустые
	if cfg.Host == "" || cfg.Port == "" {
		return nil, fmt.Errorf("SMTP_HOST and SMTP_PORT are required")
	}
	// Если пользователь и пароль не пустые, то создаем auth
	var auth smtp.Auth
	if cfg.User != "" || cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	}

	// Создаем SMTPClient
	return &SMTPClient{
		Host: cfg.Host,
		Port: cfg.Port,
		From: cfg.User,
		Auth: auth,
	}, nil
}

func (c *SMTPClient) SendMail(to []string, subject, body string) error {
	headers := []string{
		"From: " + c.From,
		"To: " + strings.Join(to, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}
	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)
	addr := net.JoinHostPort(c.Host, c.Port)
	return smtp.SendMail(addr, c.Auth, c.From, to, msg)
}
