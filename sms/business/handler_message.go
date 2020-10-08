package business

import (
	"math/rand"
	"strconv"

	"github.com/jinzhu/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

func (h *MessageHandler) CanSent(id string) (ok bool) {
	var count int
	h.db.Where(&SMSOutgoingQueue{ID: id}).Count(&count)
	return count > 0
}

func (h *MessageHandler) OnStart(slotID uint, id string, phones []string, message string) (err error) {
	tx := h.db.Begin()
	slot := new(Slot)
	if err = h.findCard(tx, slotID, slot); err != nil {
		return
	}
	cardId := slot.Card.ID
	for _, phone := range phones {
		outbox := &SMSOutbox{
			CardID:  cardId,
			Phone:   phone,
			Message: message,
			State:   SMSOutboxReady,
		}
		queue := &SMSOutgoingQueue{
			ID:        id,
			Index:     strconv.Itoa(int(rand.Int31())),
			CardID:    cardId,
			SMSOutbox: outbox,
		}
		if err := tx.Save(queue).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (h *MessageHandler) OnPassword(id string) (password string, err error) {
	queue := &SMSOutgoingQueue{ID: id}
	if err = h.db.First(queue).Error; err != nil {
		return
	}
	slot := &Slot{CardID: &queue.Card.ID}
	if err = h.db.First(slot).Error; err != nil {
		return
	}
	password = slot.Password
	return
}

func (h *MessageHandler) OnSend(id string) (index, phone string, err error) {
	// select q.index, q.sms_outbox_id, b.phone from sms_outgoing_queue as q
	// inner join sms_outbox as b on b.id = q.sms_outbox_id
	// where b.state = "READY"
	// limit 1
	panic("not impl")
}

func (h *MessageHandler) OnSendFailed(id, index, message string) (err error) {
	if message == "SENDID" {
		return
	}
	tx := h.db.Begin()
	queue := &SMSOutgoingQueue{ID: id, Index: index}
	if err = h.db.First(queue).Error; err != nil {
		return
	}
	queue.SMSOutbox.Error = &message
	queue.SMSOutbox.State = SMSOutboxFailed
	if err = tx.Delete(queue.SMSOutbox).Error; err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func (h *MessageHandler) OnSendWait(id, index string) (err error) {
	tx := h.db.Begin()
	queue := &SMSOutgoingQueue{ID: id, Index: index}
	if err = h.db.First(queue).Error; err != nil {
		return
	}
	queue.SMSOutbox.State = SMSOutboxWaiting
	if err = tx.Save(queue.SMSOutbox).Error; err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func (h *MessageHandler) OnSendSent(id, index, deliverID string) (err error) {
	tx := h.db.Begin()
	queue := &SMSOutgoingQueue{ID: id, Index: index}
	if err = h.db.First(queue).Error; err != nil {
		return
	}
	queue.SMSOutbox.State = SMSOutboxSent
	queue.DeliverID = deliverID
	if err = tx.Save(queue).Error; err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func (h *MessageHandler) OnEnd(id string) error {
	return h.db.Delete(&SMSOutgoingQueue{ID: id}).Error
}

func (h *MessageHandler) OnUSSD(slotID uint, command, message string, err error) {
	if err != nil {
		return
	}
	slot := new(Slot)
	if err = h.findCard(h.db, slotID, slot); err != nil {
		return
	}
	h.db.Save(&USSDMessage{
		CardID:  slot.Card.ID,
		Command: command,
		Message: message,
	})
}

func (h *MessageHandler) findCard(tx *gorm.DB, id uint, slot *Slot) (err error) {
	if err = h.db.First(slot, id).Error; err != nil {
		return
	}
	if slot.CardID == nil {
		return gorm.ErrRecordNotFound
	}
	slot.Card = new(Card)
	return h.db.First(slot.Card, *slot.CardID).Error
}
