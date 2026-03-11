package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Channel defines how a notification is delivered.
type Channel interface {
	Send(ctx context.Context, notifiable Notifiable, n Notification) error
}

// Notification is the interface all notifications implement.
type Notification interface {
	Channels() []string // which channels: "database", "mail"
}

// DatabaseNotification is implemented by notifications that store data
// in the database.
type DatabaseNotification interface {
	ToDatabase(notifiable Notifiable) (DatabaseMessage, error)
}

// MailNotification is implemented by notifications that send email.
type MailNotification interface {
	ToMail(notifiable Notifiable) (MailMessage, error)
}

// Notifiable is the entity receiving the notification (usually a User).
type Notifiable interface {
	NotifiableID() uint
	NotifiableEmail() string
}

// DatabaseMessage is the payload stored in the notifications table.
type DatabaseMessage struct {
	Type string
	Data map[string]interface{}
}

// MailMessage is the email content to send.
type MailMessage struct {
	Subject string
	Body    string // HTML
}

// NotificationRecord is the GORM model for stored notifications.
type NotificationRecord struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	NotifiableID uint       `gorm:"index" json:"notifiable_id"`
	Type         string     `json:"type"`
	Data         string     `gorm:"type:text" json:"data"` // JSON string
	ReadAt       *time.Time `json:"read_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (NotificationRecord) TableName() string { return "notifications" }

// Notifier dispatches notifications through registered channels.
type Notifier struct {
	channels map[string]Channel
}

// NewNotifier creates a Notifier with no channels registered.
func NewNotifier() *Notifier {
	return &Notifier{channels: make(map[string]Channel)}
}

// RegisterChannel adds a named channel to the notifier.
func (n *Notifier) RegisterChannel(name string, ch Channel) {
	n.channels[name] = ch
}

// Send dispatches a notification to the notifiable through each of the
// notification's declared channels.
func (n *Notifier) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	for _, name := range notification.Channels() {
		ch, ok := n.channels[name]
		if !ok {
			return fmt.Errorf("notification: channel %q not registered", name)
		}
		if err := ch.Send(ctx, notifiable, notification); err != nil {
			return fmt.Errorf("notification: channel %q failed: %w", name, err)
		}
	}
	return nil
}

// DatabaseChannel stores notifications in the database via GORM.
type DatabaseChannel struct {
	db *gorm.DB
}

// NewDatabaseChannel creates a DatabaseChannel.
func NewDatabaseChannel(db *gorm.DB) *DatabaseChannel {
	return &DatabaseChannel{db: db}
}

// Send persists the notification to the notifications table.
func (dc *DatabaseChannel) Send(ctx context.Context, notifiable Notifiable, n Notification) error {
	dn, ok := n.(DatabaseNotification)
	if !ok {
		return fmt.Errorf("notification does not implement DatabaseNotification")
	}
	msg, err := dn.ToDatabase(notifiable)
	if err != nil {
		return err
	}
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return fmt.Errorf("notification: failed to marshal data: %w", err)
	}

	record := NotificationRecord{
		NotifiableID: notifiable.NotifiableID(),
		Type:         msg.Type,
		Data:         string(data),
		CreatedAt:    time.Now(),
	}
	return dc.db.WithContext(ctx).Create(&record).Error
}

// MailSender is a minimal interface for sending email, matching mail.Mailer.Send.
type MailSender interface {
	Send(to, subject, htmlBody string) error
}

// MailChannel sends notifications via email.
type MailChannel struct {
	mailer MailSender
}

// NewMailChannel creates a MailChannel with the given mailer.
func NewMailChannel(mailer MailSender) *MailChannel {
	return &MailChannel{mailer: mailer}
}

// Send delivers the notification via email.
func (mc *MailChannel) Send(ctx context.Context, notifiable Notifiable, n Notification) error {
	mn, ok := n.(MailNotification)
	if !ok {
		return fmt.Errorf("notification does not implement MailNotification")
	}
	msg, err := mn.ToMail(notifiable)
	if err != nil {
		return err
	}
	return mc.mailer.Send(notifiable.NotifiableEmail(), msg.Subject, msg.Body)
}
