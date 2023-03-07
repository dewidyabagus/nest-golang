package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	From   string
	dialer *gomail.Dialer
}

func NewMailer(cfg MailConfig) *mailer {
	return &mailer{
		From:   cfg.SenderName,
		dialer: gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.Email, cfg.Password),
	}
}

func (m *mailer) SendLoginNotify(info UserNotify) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", info.Email)
	msg.SetHeader("Subject", "Notifikasi Masuk Dashboard")
	msg.SetBody("text/html", fmt.Sprintf(
		`Hello, %s %s <br>
		Account kamu baru saja login .....
		`,
		info.FirstName, info.LastName))

	return m.dialer.DialAndSend(msg)
}

func (m *mailer) SendAccountActivation(info UserNotify) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", info.Email)
	msg.SetHeader("Subject", "Notifikasi Aktivasi Akun")
	msg.SetBody("text/html", fmt.Sprintf(
		`Hello, %s %s <br>
		Aktifkan akun kamu dengan klik link berikut .....
		`,
		info.FirstName, info.LastName))

	return m.dialer.DialAndSend(msg)
}

func (m *mailer) Ping() error {
	closer, err := m.dialer.Dial()
	if err != nil {
		return err
	}

	return closer.Close()
}
