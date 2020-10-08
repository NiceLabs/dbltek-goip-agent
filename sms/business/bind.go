package business

import (
	"github.com/jinzhu/gorm"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

func BindDatabase(db *gorm.DB, serve *udpstack.Server) {
	db.AutoMigrate(
		new(Slot),
		new(Card),
		new(Record),
		new(SMSInbox),
		new(SMSOutbox),
		new(SMSOutgoingQueue),
		new(USSDMessage),
		new(CronTask),
	)
	notifyHandler := &NotifyHandler{db: db, serve: serve}
	messageHandler := &MessageHandler{db: db}
	serve.OnNotify = notifyHandler
	serve.OnBulkSMS = messageHandler
	serve.OnUSSD = messageHandler
}

func BindCommand(db *gorm.DB, serve *udpstack.Server, id string) (cmd *udpstack.Command, err error) {
	slot := new(Slot)
	err = db.
		Where(Slot{DeviceID: id}).
		Select([]string{"address", "password", "id"}).
		First(slot).
		Error
	if err != nil {
		return
	}
	return slot.Command(serve)
}
