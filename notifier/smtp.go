package notifier

import (
	"fmt"
	"net/smtp"
)

// sendMailFunc is a package-level variable to allow mocking in tests
var sendMailFunc = smtp.SendMail

type SMTPNotifier struct {
	host        string
	port        string
	user        string
	pass        string
	targetEmail string
	fromEmail   string
}

// NewSMTPNotifier creates a new client for sending SMS via Email gateway
func NewSMTPNotifier(host, port, user, pass, targetEmail, fromEmail string) *SMTPNotifier {
	return &SMTPNotifier{
		host:        host,
		port:        port,
		user:        user,
		pass:        pass,
		targetEmail: targetEmail,
		fromEmail:   fromEmail,
	}
}

// SendSMS implements the server.Notifier interface
func (s *SMTPNotifier) SendSMS(message string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	
	// Set up authentication information if provided
	var auth smtp.Auth
	if s.user != "" && s.pass != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	// Construct RFC 822 formatted email
	body := fmt.Sprintf("To: %s\r\nSubject: BART Alert\r\n\r\n%s", s.targetEmail, message)

	err := sendMailFunc(addr, auth, s.fromEmail, []string{s.targetEmail}, []byte(body))
	if err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}
	
	return nil
}
