package exchangesmtp

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Mail is a struct of plain text email.
type Mail struct {
	mime        string
	From        string
	To          []string
	Subject     string
	Body        string
	ContentType string
	Attachment  Attachment
}

// AttachmentFile representing an attachment in mail
type AttachmentFile struct {
	Name        string
	ContentType string
	Body        []byte
}

// Attachment is a wrapper for AttachmentFile
type Attachment struct {
	File     []AttachmentFile
	WithFile bool
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

	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = mail.From
	Header["To"] = strings.Join(mail.To, ";")
	Header["Cc"] = strings.Join([]string{}, ";")
	Header["Bcc"] = strings.Join([]string{}, ";")
	Header["Subject"] = mail.Subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + mail.ContentType + "\r\n"
	body += "\r\n" + mail.Body + "\r\n"
	buffer.WriteString(body)

	if mail.Attachment.WithFile {
		for _, file := range mail.Attachment.File {
			attachment := "\r\n--" + boundary + "\r\n"
			attachment += "Content-Transfer-Encoding:base64\r\n"
			attachment += "Content-Disposition:attachment\r\n"
			attachment += "Content-Type:" + file.ContentType + ";name=\"" + file.Name + "\"\r\n"
			_, err := buffer.WriteString(attachment)
			if err != nil {
				return err
			}

			if len(file.Body) > 0 {
				if err := writeBytes(buffer, file.Body); err != nil {
					return err
				}
			} else {
				if err := writeFile(buffer, file.Name); err != nil {
					return err
				}
			}
		}
	}
	buffer.WriteString("\r\n--" + boundary + "--")

	if err := smtp.SendMail(m.server, m.auth, mail.From, mail.To, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)

	return header
}

func writeFile(buffer *bytes.Buffer, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err = writeBytes(buffer, file); err != nil {
		return err
	}

	return nil
}

func writeBytes(buffer *bytes.Buffer, file []byte) error {
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}

	return nil
}
