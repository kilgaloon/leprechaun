package notifications

import (
	"log"
	"net/smtp"
	"net/url"
)

// Email is type of notification
type Email struct {
	message           string
	title             string
	host              string
	username          string
	password          string
	notificationEmail string
}

// EmailConfig defines interface that satisfies NewEmail
type EmailConfig interface {
	GetSMTPHost() string
	GetSMTPUsername() string
	GetSMTPPassword() string
	GetNotificationsEmail() string
}

// SetMessage will be used as body for email
func (e *Email) SetMessage(m string) {
	e.message = m
}

// GetMessage get message of notification
func (e Email) GetMessage() string {
	return e.message
}

// SetTitle will be used as title for email
func (e *Email) SetTitle(t string) {
	e.title = t
}

// GetTitle get message of notification
func (e Email) GetTitle() string {
	return e.title
}

// Send sends notification
func (e Email) Send() error {
	to := e.notificationEmail
	u, err := url.Parse("http://" + e.host)
	if err != nil {
		panic(err)
	}

	msg := "From: " + e.username + "\n" +
		"To: " + to + "\n" +
		"Subject: " + e.GetTitle() + "\n\n" +
		e.GetMessage()

	err = smtp.SendMail(e.host,
		smtp.PlainAuth("", e.username, e.password, u.Hostname()),
		e.username, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}

// NewEmail create new email notification
// params from this function is used to build smtp.Client
func NewEmail(cfg EmailConfig) *Email {
	e := &Email{
		host:              cfg.GetSMTPHost(),
		username:          cfg.GetSMTPUsername(),
		password:          cfg.GetSMTPPassword(),
		notificationEmail: cfg.GetNotificationsEmail(),
	}

	return e
}
