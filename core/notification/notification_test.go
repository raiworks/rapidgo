package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// ─── Test Doubles ─────────────────────────────────────────────────────────────

type testUser struct {
	id    uint
	email string
}

func (u testUser) NotifiableID() uint      { return u.id }
func (u testUser) NotifiableEmail() string { return u.email }

type testNotification struct {
	channels []string
	dbMsg    DatabaseMessage
	mailMsg  MailMessage
}

func (n testNotification) Channels() []string { return n.channels }
func (n testNotification) ToDatabase(_ Notifiable) (DatabaseMessage, error) {
	return n.dbMsg, nil
}
func (n testNotification) ToMail(_ Notifiable) (MailMessage, error) {
	return n.mailMsg, nil
}

type mockMailer struct {
	sent []struct{ To, Subject, Body string }
}

func (m *mockMailer) Send(to, subject, body string) error {
	m.sent = append(m.sent, struct{ To, Subject, Body string }{to, subject, body})
	return nil
}

type failingMailer struct{}

func (failingMailer) Send(_, _, _ string) error { return fmt.Errorf("smtp down") }

func setupNotificationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&NotificationRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestNewNotifier(t *testing.T) {
	n := NewNotifier()
	if n == nil {
		t.Fatal("NewNotifier returned nil")
	}
}

// T-024: Database channel persists notification
func TestDatabaseChannel_Send(t *testing.T) {
	db := setupNotificationDB(t)
	ch := NewDatabaseChannel(db)

	user := testUser{id: 42, email: "user@test.com"}
	notif := testNotification{
		channels: []string{"database"},
		dbMsg: DatabaseMessage{
			Type: "new_follower",
			Data: map[string]interface{}{"follower_id": float64(7), "name": "Alice"},
		},
	}

	err := ch.Send(context.Background(), user, notif)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	var record NotificationRecord
	db.First(&record)
	if record.NotifiableID != 42 {
		t.Errorf("NotifiableID = %d, want 42", record.NotifiableID)
	}
	if record.Type != "new_follower" {
		t.Errorf("Type = %q, want 'new_follower'", record.Type)
	}
	var data map[string]interface{}
	json.Unmarshal([]byte(record.Data), &data)
	if data["name"] != "Alice" {
		t.Errorf("Data.name = %v, want Alice", data["name"])
	}
}

// T-025: Mail channel sends email
func TestMailChannel_Send(t *testing.T) {
	mailer := &mockMailer{}
	ch := NewMailChannel(mailer)

	user := testUser{id: 1, email: "user@test.com"}
	notif := testNotification{
		channels: []string{"mail"},
		mailMsg:  MailMessage{Subject: "Welcome!", Body: "<p>Hello</p>"},
	}

	err := ch.Send(context.Background(), user, notif)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}
	if len(mailer.sent) != 1 {
		t.Fatalf("expected 1 sent email, got %d", len(mailer.sent))
	}
	if mailer.sent[0].To != "user@test.com" {
		t.Errorf("To = %q, want user@test.com", mailer.sent[0].To)
	}
	if mailer.sent[0].Subject != "Welcome!" {
		t.Errorf("Subject = %q, want Welcome!", mailer.sent[0].Subject)
	}
}

// T-026: Notifier dispatches to multiple channels
func TestNotifier_MultiChannel(t *testing.T) {
	db := setupNotificationDB(t)
	mailer := &mockMailer{}

	notifier := NewNotifier()
	notifier.RegisterChannel("database", NewDatabaseChannel(db))
	notifier.RegisterChannel("mail", NewMailChannel(mailer))

	user := testUser{id: 5, email: "multi@test.com"}
	notif := testNotification{
		channels: []string{"database", "mail"},
		dbMsg: DatabaseMessage{
			Type: "welcome",
			Data: map[string]interface{}{"msg": "hi"},
		},
		mailMsg: MailMessage{Subject: "Welcome!", Body: "<p>Hi</p>"},
	}

	err := notifier.Send(context.Background(), user, notif)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	// Check database
	var count int64
	db.Model(&NotificationRecord{}).Count(&count)
	if count != 1 {
		t.Errorf("db records = %d, want 1", count)
	}

	// Check mail
	if len(mailer.sent) != 1 {
		t.Errorf("emails sent = %d, want 1", len(mailer.sent))
	}
}

// T-027: Unregistered channel returns error
func TestNotifier_UnregisteredChannel(t *testing.T) {
	notifier := NewNotifier()

	user := testUser{id: 1, email: "x@test.com"}
	notif := testNotification{channels: []string{"push"}}

	err := notifier.Send(context.Background(), user, notif)
	if err == nil {
		t.Fatal("expected error for unregistered channel")
	}
}

// T-028: Channel Send error propagates
func TestNotifier_ChannelError(t *testing.T) {
	notifier := NewNotifier()
	notifier.RegisterChannel("mail", NewMailChannel(failingMailer{}))

	user := testUser{id: 1, email: "x@test.com"}
	notif := testNotification{
		channels: []string{"mail"},
		mailMsg:  MailMessage{Subject: "Test", Body: "body"},
	}

	err := notifier.Send(context.Background(), user, notif)
	if err == nil {
		t.Fatal("expected error from failing mailer")
	}
}
