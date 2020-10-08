package business

import (
	"net"

	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

func (h *NotifyHandler) OnBCCHUpdate(addr net.Addr, u *udpstack.BCCHUpdate) (err error) {
	return h.OnCellsUpdate(addr, &udpstack.CellsUpdate{
		RequestID: u.RequestID,
		ID:        u.ID,
		Password:  u.Password,
		Cells:     u.Cells,
	})
}

func (h *NotifyHandler) OnCellsUpdate(addr net.Addr, u *udpstack.CellsUpdate) (err error) {
	slot := new(Slot)
	result := h.db.First(slot, Slot{DeviceID: u.ID})
	if err := result.Error; err != nil {
		return err
	}
	slot.Address = addr.String()
	slot.Password = u.Password
	slot.Metadata.Cells = u.Cells
	return
}

func (h *NotifyHandler) OnATUpdate(net.Addr, *udpstack.ATUpdate) error {
	return nil
}
