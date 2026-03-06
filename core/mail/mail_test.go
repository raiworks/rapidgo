package mail

import (
	"testing"
)

// TC-01: NewMailer reads env vars
func TestNewMailer_ReadsEnv(t *testing.T) {
	t.Setenv("MAIL_HOST", "smtp.test.com")
	t.Setenv("MAIL_PORT", "465")
	t.Setenv("MAIL_USERNAME", "user@test.com")
	t.Setenv("MAIL_PASSWORD", "s3cret")
	t.Setenv("MAIL_FROM_ADDRESS", "noreply@test.com")
	t.Setenv("MAIL_FROM_NAME", "TestApp")

	m := NewMailer()

	if m.Host != "smtp.test.com" {
		t.Fatalf("Host: expected 'smtp.test.com', got %q", m.Host)
	}
	if m.Port != 465 {
		t.Fatalf("Port: expected 465, got %d", m.Port)
	}
	if m.Username != "user@test.com" {
		t.Fatalf("Username: expected 'user@test.com', got %q", m.Username)
	}
	if m.Password != "s3cret" {
		t.Fatalf("Password: expected 's3cret', got %q", m.Password)
	}
	if m.From != "noreply@test.com" {
		t.Fatalf("From: expected 'noreply@test.com', got %q", m.From)
	}
	if m.FromName != "TestApp" {
		t.Fatalf("FromName: expected 'TestApp', got %q", m.FromName)
	}
}

// TC-02: NewMailer defaults port to 587 when unset
func TestNewMailer_DefaultPort(t *testing.T) {
	t.Setenv("MAIL_PORT", "")

	m := NewMailer()
	if m.Port != 587 {
		t.Fatalf("expected default port 587, got %d", m.Port)
	}
}

// TC-03: NewMailer parses port from string
func TestNewMailer_ParsesPort(t *testing.T) {
	t.Setenv("MAIL_PORT", "2525")

	m := NewMailer()
	if m.Port != 2525 {
		t.Fatalf("expected port 2525, got %d", m.Port)
	}
}

// TC-04: NewMailer handles invalid port gracefully
func TestNewMailer_InvalidPortFallback(t *testing.T) {
	t.Setenv("MAIL_PORT", "abc")

	m := NewMailer()
	if m.Port != 587 {
		t.Fatalf("expected fallback port 587 for invalid input, got %d", m.Port)
	}
}

// TC-05: Send returns error with invalid host
func TestMailer_Send_InvalidHost(t *testing.T) {
	m := &Mailer{
		Host:     "invalid.host.test",
		Port:     587,
		Username: "user",
		Password: "pass",
		From:     "test@example.com",
		FromName: "Test",
	}

	err := m.Send("recipient@example.com", "Test Subject", "<p>Hello</p>")
	if err == nil {
		t.Fatal("expected error when sending to invalid host")
	}
}
