package exchangesmtp

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// Mail is a struct of plain text email.
type Mail struct {
	mime    string
	From    string
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func (m *Mail) String() string {
	return fmt.Sprintf("%sFrom: %s\nTo: %s\nSubject: %s\n\n%s", m.mime, m.From, strings.Join(m.To, ","), m.Subject, m.Body)
}

// MailSender uses for send plain text emails.
type MailSender struct {
	auth   smtp.Auth
	server string
}

// NewMailSender is constructor for MailSender.
func NewMailSender(username, password, server string) *MailSender {
	return &MailSender{LoginAuth(username, password), server}
}

// SendToList is send simple plain text email to list recipients.
func (m *MailSender) SendToList(mail Mail) error {
	if len(mail.To) == 0 {
		return errors.New("recipient list is empty")
	}
	if len(strings.TrimSpace(mail.Body)) == 0 {
		return errors.New("message is empty")
	}
	addMime(&mail)

	return smtp.SendMail(m.server, m.auth, mail.From, mail.To, []byte(mail.String()))
}

func addMime(mail *Mail) {
	if mail.IsHTML {
		mail.mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	}
}
