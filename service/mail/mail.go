package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

const (
	mime    = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: Your Shorters' OTP!\n"
	host    = "smtp.gmail.com"
	server  = "smtp.gmail.com:587"
)

type MailService interface {
	SendOTP(mail string, otp string) error
}

type mailService struct {
	mail string
	pass string
}

func NewMailService() MailService {
	return &mailService{
		mail: os.Getenv("mail"),
		pass: os.Getenv("mail_pass"),
	}
}

func (s *mailService) SendOTP(mail string, otp string) error {
	t, err := template.ParseFiles("view/mail.html")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, struct{ OTP string }{otp})
	if err != nil {
		return err
	}
	content := buf.String()
	body := []byte(subject + mime + content)
	auth := smtp.PlainAuth("", s.mail, s.pass, host)
	err = smtp.SendMail(server, auth, s.mail, []string{mail}, body)
	if err != nil {
		return fmt.Errorf("error from SMTP server: %v", err)
	}
	return nil
}
