package business

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"

	"golang.org/x/xerrors"
)

//noinspection SpellCheckingInspection
type SlotVoIP struct {
	Status     string `json:",omitempty"`
	State      string `json:",omitempty"`
	TotalTime  *int   `json:",omitempty"`
	RemainTime *int   `json:",omitempty"`
	IdleTime   int    `json:",omitempty"`
	RejectCall bool   `json:",omitempty"`
}

func (c SlotVoIP) Unacceptable() bool {
	return c.RejectCall ||
		c.Status == udpstack.VoIPStatusLogout ||
		c.Status == udpstack.VoIPStatusDown
}

func (c SlotVoIP) Value() (driver.Value, error) {
	encoded, err := json.Marshal(c)
	return string(encoded), err
}

func (c *SlotVoIP) Scan(value interface{}) error {
	if input, ok := value.(string); ok {
		return json.Unmarshal([]byte(input), c)
	}
	return xerrors.Errorf("Failed to unmarshal value: %#v", value)
}
