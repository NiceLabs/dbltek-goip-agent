package business

import (
	"net"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

//noinspection SpellCheckingInspection
type Slot struct {
	gorm.Model
	Address  string
	DeviceID string `gorm:"unique"`
	Password string
	CardID   *uint
	Card     *Card
	VoIP     SlotVoIP     `gorm:"type:json1;column:voip"`
	Metadata SlotMetadata `gorm:"type:json1"`
	Disabled bool
}

func (s *Slot) Outdated() bool {
	return time.Now().Sub(s.UpdatedAt).Minutes() >= 2
}

func (s *Slot) Command(server *udpstack.Server) (command *udpstack.Command, err error) {
	addr, err := net.ResolveUDPAddr("udp", s.Address)
	if err != nil {
		return
	}
	command = server.Command(addr, s.ID, s.Password)
	return
}
