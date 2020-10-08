package business

import (
	"database/sql/driver"
	"encoding/json"

	. "github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"

	"golang.org/x/xerrors"
)

//noinspection SpellCheckingInspection
type SlotMetadata struct {
	RSSI            int       `json:",omitempty"` // Cell: RSSI
	Status          string    `json:",omitempty"` // Cell: Registered Status
	Carrier         string    `json:",omitempty"` // Cell: Registered Carrier
	CGATT           bool      `json:",omitempty"` // Cell: PS attach or detach
	Cells           []string  `json:",omitempty"` // Cell: Cell-list
	CellInfo        *CellInfo `json:",omitempty"` // Cell: Current Cell info
	IMEI            string    `json:",omitempty"` // Slot: IMEI
	OutcallInterval *int      `json:",omitempty"` // GoIP: Outcall interval
}

func (c SlotMetadata) Value() (driver.Value, error) {
	encoded, err := json.Marshal(c)
	return string(encoded), err
}

func (c *SlotMetadata) Scan(value interface{}) error {
	if input, ok := value.(string); ok {
		return json.Unmarshal([]byte(input), c)
	}
	return xerrors.Errorf("Failed to unmarshal value: %#v", value)
}
