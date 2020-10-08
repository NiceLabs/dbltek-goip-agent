package business

type SMSOutgoingQueue struct {
	ID          string `gorm:"primary_key"`
	Index       string `gorm:"primary_key"`
	CardID      uint
	Card        Card
	SMSOutboxID uint
	SMSOutbox   *SMSOutbox
	DeliverID   string
}

func (q *SMSOutgoingQueue) TableName() string {
	return "sms_outgoing_queue"
}
