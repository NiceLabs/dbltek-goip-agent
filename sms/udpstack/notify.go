package udpstack

//RegistrationUpdate
//noinspection SpellCheckingInspection
type RegistrationUpdate struct {
	RequestID      string     `field:"req"`
	ID             string     `field:"id"`
	Password       string     `field:"pass"`
	Phone          string     `field:"num"`            // Phone number
	Signal         int        `field:"signal"`         //  GSM Signal: Signal strength, AT+CSQ format
	Status         VoIPStatus `field:"gsm_status"`     //  GSM Status: LOGIN, LOGOUT
	VoIPStatus     VoIPStatus `field:"voip_status"`    // VoIP Status: LOGIN, LOGOUT, UP, DOWN
	VoIPState      string     `field:"voip_state"`     // VoIP State
	VoIPRemainTime *int       `field:"remain_time"`    // VoIP Remain time (minutes), with calling update
	VoIPIdleTime   int        `field:"idle"`           // VoIP Idle time (minutes), with calling reset
	IMEI           string     `field:"imei"`           // IMEI in Slot
	IMSI           string     `field:"imsi"`           // IMSI in Card
	ICCID          string     `field:"iccid"`          // ICCID in Card
	Carrier        string     `field:"pro"`            // Carrier in Slot, (e.q: CHINA MOBILE)
	SMBLogin       bool       `field:"SMB_LOGIN"`      // Registered to SIM Server
	SMSLogin       bool       `field:"SMS_LOGIN"`      // Registered to SMS Server
	CellInfo       string     `field:"CELLINFO"`       // Format: "LAC:5815,CELL ID:9487"
	CGATT          bool       `field:"CGATT"`          // PS attach or detach
	Disabled       bool       `field:"disable_status"` // Disabled Channel
}

//noinspection SpellCheckingInspection
func (r *RegistrationUpdate) RSSI() int {
	if r.Signal == 99 {
		return 0
	}
	return r.Signal*2 - 113
}
