package business

import (
	"github.com/jinzhu/gorm"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

type Record struct {
	gorm.Model
	CardID      uint
	Card        Card
	Direction   udpstack.Direction
	Phone       string `gorm:"index"`
	CallTime    *int
	HangupCause *string
}

func (r *Record) SetCause(cause string) {
	r.HangupCause = &cause
}
