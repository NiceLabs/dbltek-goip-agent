package business

import (
	"net"
	"reflect"

	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"

	"github.com/jinzhu/gorm"
)

type NotifyHandler struct {
	db    *gorm.DB
	serve *udpstack.Server
}

func (h *NotifyHandler) OnNotify(addr net.Addr, update interface{}) error {
	name := "On" + reflect.TypeOf(update).Elem().Name()
	out := reflect.ValueOf(h).MethodByName(name).Call([]reflect.Value{
		reflect.ValueOf(addr),
		reflect.ValueOf(update),
	})
	if err := out[0]; !err.IsNil() {
		return err.Interface().(error)
	}
	return nil
}

func (h *NotifyHandler) OnRegistrationUpdate(addr net.Addr, u *udpstack.RegistrationUpdate) (err error) {
	slot := new(Slot)
	if err = h.db.FirstOrCreate(slot, Slot{DeviceID: u.ID}).Error; err != nil {
		return
	}
	if u.IMSI != "" {
		slot.Card = new(Card)
		if err = h.db.FirstOrCreate(slot.Card, Card{IMSI: u.IMSI}).Error; err != nil {
			return
		}
		slot.Card.Phone = formatPhoneNumber(u.IMSI, u.Phone)
		slot.Card.ICCID = u.ICCID
	}
	slot.Address = addr.String()
	slot.Password = u.Password
	slot.VoIP.Status = string(u.VoIPStatus)
	slot.VoIP.State = u.VoIPState
	slot.VoIP.RemainTime = u.VoIPRemainTime
	slot.VoIP.IdleTime = u.VoIPIdleTime
	slot.Metadata.RSSI = u.RSSI()
	slot.Metadata.Status = string(u.Status)
	slot.Metadata.Carrier = u.Carrier
	slot.Metadata.IMEI = u.IMEI
	slot.Metadata.CGATT = u.CGATT
	slot.Metadata.CellInfo = udpstack.ParseCellInfo(u.CellInfo)
	slot.Disabled = u.Disabled
	return h.db.Save(slot).Error
}

func (h *NotifyHandler) findCard(id string, slot *Slot) (err error) {
	if err = h.db.First(slot, Slot{DeviceID: id}).Error; err != nil {
		return
	}
	if slot.CardID == nil {
		return gorm.ErrRecordNotFound
	}
	slot.Card = new(Card)
	return h.db.First(slot.Card, *slot.CardID).Error
}
