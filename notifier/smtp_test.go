package notifier

import (
	"net/smtp"
	"strings"
	"testing"
)

func TestSendSMS_Success(t *testing.T) {
	var capturedAddr string
	var capturedFrom string
	var capturedTo []string
	var capturedMsg []byte

	// Mock the sendMail variable to intercept the call
	originalSendMail := sendMailFunc
	defer func() { sendMailFunc = originalSendMail }()

	sendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedAddr = addr
		capturedFrom = from
		capturedTo = to
		capturedMsg = msg
		return nil
	}

	notifier := NewSMTPNotifier("smtp.example.com", "587", "user", "pass", "target@vtext.com", "sender@example.com")
	
	err := notifier.SendSMS("Test warning message")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if capturedAddr != "smtp.example.com:587" {
		t.Errorf("Expected addr smtp.example.com:587, got %s", capturedAddr)
	}
	if capturedFrom != "sender@example.com" {
		t.Errorf("Expected from sender@example.com, got %s", capturedFrom)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "target@vtext.com" {
		t.Errorf("Expected to [target@vtext.com], got %v", capturedTo)
	}

	msgStr := string(capturedMsg)
	if !strings.Contains(msgStr, "To: target@vtext.com") {
		t.Errorf("Message missing To header: %s", msgStr)
	}
	if !strings.Contains(msgStr, "Subject: BART Alert") {
		t.Errorf("Message missing Subject header: %s", msgStr)
	}
	if !strings.Contains(msgStr, "\r\n\r\nTest warning message") {
		t.Errorf("Message missing body: %s", msgStr)
	}
}
