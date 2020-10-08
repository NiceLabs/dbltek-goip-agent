package business

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SMSInbox struct {
	gorm.Model
	CardID  uint
	Card    Card
	Phone   string
	Message string
}

func (b *SMSInbox) TableName() string {
	return "sms_inbox"
}

type SMSOutboxState string

const (
	SMSOutboxReady     SMSOutboxState = "READY"
	SMSOutboxSending   SMSOutboxState = "SENDING"
	SMSOutboxWaiting   SMSOutboxState = "WAITING"
	SMSOutboxSent      SMSOutboxState = "SENT"
	SMSOutboxFailed    SMSOutboxState = "FAILED"
	SMSOutboxDelivered SMSOutboxState = "DELIVERED"
)

type SMSOutbox struct {
	gorm.Model
	CardID   uint
	Card     Card
	Phone    string
	Message  string
	SendTime *time.Time
	Error    *string
	State    SMSOutboxState
}

func (b *SMSOutbox) TableName() string {
	return "sms_outbox"
}

type USSDMessage struct {
	gorm.Model
	CardID  uint
	Card    Card
	Command string
	Message string
}
