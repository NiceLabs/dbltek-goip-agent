package udpstack

import (
	"regexp"
)

type (
	// State changed
	//   Call-in: IDLE -> INCOMING:\d+ -> ACTIVE -> IDLE
	//   Call-out accepted: IDLE -> DIALING:\d+ -> ALERTING:\d+ -> CONNECTED:\d+ -> IDLE
	//   Call-out rejected: IDLE -> DIALING:\d+ -> ALERTING:\d+ -> IDLE
	//   BUSY?
	StateUpdate struct {
		RequestID string `field:"STATE"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		State     string `field:"gsm_remain_state"`
	}
	RecordUpdate struct {
		RequestID string    `field:"RECORD"`
		ID        string    `field:"id"`
		Password  string    `field:"password"`
		Direction Direction `field:"dir"`
		Phone     string    `field:"num"`
	}
	HangupUpdate struct {
		RequestID string `field:"HANGUP"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Return    string `field:"num"`
	}
	CallTimeUpdate struct {
		RequestID string `field:"EXPIRY"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Time      *int   `field:"exp"`
	}
	RemainTimeUpdate struct {
		RequestID string `field:"REMAIN"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Time      *int   `field:"gsm_remain_time"`
	}
	GSMUpdate struct {
		RequestID       string  `field:"GSM"`
		Phone           *string `field:"get_gsm_num"`
		VoIPState       string  `field:"get_gsm_state"`
		VoIPTotalTime   *int    `field:"gsm_exp_time"`
		VoIPRemainTime  *int    `field:"gsm_remain_time"`
		OutcallInterval *int    `field:"out_call_interval"`
		IMEI            string  `field:"imei"`
		Disabled        bool    `field:"module_down"`
	}
)

func (u *HangupUpdate) ParseCause() (cause string) {
	re := regexp.MustCompile(`,cause:(.*)?`)
	matched := re.FindStringSubmatch(u.Return)
	cause = matched[1]
	return
}
