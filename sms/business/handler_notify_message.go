package business

import (
	"net"

	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

func (h *NotifyHandler) OnReceiveSMSUpdate(addr net.Addr, u *udpstack.ReceiveSMSUpdate) (err error) {
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	message := &SMSInbox{
		CardID:  slot.Card.ID,
		Phone:   formatPhoneNumber(slot.Card.IMSI, u.Phone),
		Message: u.Message,
	}
	return h.db.Save(message).Error
}

func (h *NotifyHandler) OnDeliverSMSUpdate(addr net.Addr, u *udpstack.DeliverSMSUpdate) (err error) {
	if u.State != 0 {
		return
	}
	tx := h.db.Begin()
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	queue := &SMSOutgoingQueue{
		CardID:    slot.Card.ID,
		DeliverID: u.SMSNo,
	}
	if err = tx.First(&queue).Error; err != nil {
		return
	}
	queue.SMSOutbox.State = SMSOutboxDelivered

	if err = tx.Save(queue.SMSOutbox).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Delete(queue).Error; err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func (h *NotifyHandler) OnUSSNUpdate(addr net.Addr, u *udpstack.USSNUpdate) (err error) {
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	message := &USSDMessage{
		CardID:  slot.Card.ID,
		Message: u.Message,
	}
	return h.db.Save(message).Error
}
