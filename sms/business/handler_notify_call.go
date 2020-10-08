package business

import (
	"net"
	"time"

	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

func (h *NotifyHandler) OnRecordUpdate(addr net.Addr, u *udpstack.RecordUpdate) (err error) {
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	record := &Record{
		CardID:    slot.Card.ID,
		Direction: u.Direction,
		Phone:     formatPhoneNumber(slot.Card.IMSI, u.Phone),
	}
	return h.db.
		FirstOrCreate(
			record,
			"created_at between ? and ?",
			time.Now().Add(-time.Second),
			time.Now().Add(time.Second),
		).
		Error
}

func (h *NotifyHandler) OnStateUpdate(addr net.Addr, u *udpstack.StateUpdate) (err error) {
	slot := new(Slot)
	if err = h.db.First(slot, Slot{DeviceID: u.ID}).Error; err != nil {
		return err
	}
	if u.State != "IDLE" && slot.VoIP.Unacceptable() {
		command := h.serve.Command(addr, slot.ID, u.Password)
		//noinspection GoUnhandledErrorResult
		go command.EndCall()
	}
	slot.Address = addr.String()
	slot.Password = u.Password
	slot.VoIP.State = u.State
	return h.db.Save(slot).Error
}

func (h *NotifyHandler) OnHangupUpdate(addr net.Addr, u *udpstack.HangupUpdate) (err error) {
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	record := new(Record)
	if err = h.db.Where("card_id = ?", slot.Card.ID).Last(record).Error; err != nil {
		return err
	}
	record.SetCause(u.ParseCause())
	return h.db.Save(record).Error
}

func (h *NotifyHandler) OnCallTimeUpdate(addr net.Addr, u *udpstack.CallTimeUpdate) (err error) {
	slot := new(Slot)
	if err = h.findCard(u.ID, slot); err != nil {
		return
	}
	record := new(Record)
	if err = h.db.Where("card_id = ?", slot.Card.ID).Last(record).Error; err != nil {
		return err
	}
	record.CallTime = u.Time
	return h.db.Save(record).Error
}

func (h *NotifyHandler) OnRemainTimeUpdate(addr net.Addr, u *udpstack.RemainTimeUpdate) (err error) {
	slot := new(Slot)
	if err = h.db.First(slot, Slot{DeviceID: u.ID}).Error; err != nil {
		return err
	}
	slot.Address = addr.String()
	slot.Password = u.Password
	slot.VoIP.RemainTime = u.Time
	return h.db.Save(slot).Error
}

func (h *NotifyHandler) OnGSMUpdate(addr net.Addr, u *udpstack.GSMUpdate) (err error) {
	slot := new(Slot)
	if err = h.db.First(slot, Slot{Address: addr.String()}).Error; err != nil {
		return
	}
	slot.Address = addr.String()
	slot.VoIP.State = u.VoIPState
	slot.VoIP.TotalTime = u.VoIPTotalTime
	slot.VoIP.RemainTime = u.VoIPRemainTime
	slot.Metadata.IMEI = u.IMEI
	slot.Disabled = u.Disabled
	slot.Metadata.OutcallInterval = u.OutcallInterval
	return h.db.Save(slot).Error
}
