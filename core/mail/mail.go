package mail

import (
	"os"
	"strconv"

	gomail "github.com/wneessen/go-mail"
)

// Mailer holds SMTP configuration for sending emails.
type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// NewMailer creates a Mailer from MAIL_* environment variables.
func NewMailer() *Mailer {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil || port <= 0 {
		port = 587
	}
	return &Mailer{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		From:     os.Getenv("MAIL_FROM_ADDRESS"),
		FromName: os.Getenv("MAIL_FROM_NAME"),
	}
}

// Send delivers an HTML email to the specified recipient.
func (m *Mailer) Send(to, subject, htmlBody string) error {
	msg := gomail.NewMsg()
	if err := msg.FromFormat(m.FromName, m.From); err != nil {
		return err
	}
	if err := msg.To(to); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyString(gomail.TypeTextHTML, htmlBody)

	client, err := gomail.NewClient(m.Host,
		gomail.WithPort(m.Port),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(m.Username),
		gomail.WithPassword(m.Password),
		gomail.WithTLSPolicy(gomail.TLSMandatory),
	)
	if err != nil {
		return err
	}
	return client.DialAndSend(msg)
}
