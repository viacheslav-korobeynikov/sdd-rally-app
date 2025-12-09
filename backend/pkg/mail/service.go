package mail

import (
	"fmt"
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

func NewSMTPClient(auth smtp.Auth) *SMTPClient {
	SMTPConfig := config.NewSMTPConfig()
	return &SMTPClient{
		Host: SMTPConfig.Host,
		Port: SMTPConfig.Port,
		From: SMTPConfig.User,
		Auth: auth,
	}
}

func (c *SMTPClient) SendMail(to []string, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", strings.Join(to, ", "), subject, body))
	addr := c.Host + ":" + c.Port
	return smtp.SendMail(addr, c.Auth, c.From, to, msg)
}
