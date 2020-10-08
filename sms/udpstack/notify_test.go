package udpstack

import (
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NotifyTestCase struct {
	Packet   string
	Expected interface{}
}

var casePhone = "+861380013800"

//noinspection SpellCheckingInspection
var testCases = []NotifyTestCase{
	{
		"req:3;id:goip-1;pass:goip;num:;signal:;gsm_status:LOGIN;voip_status:LOGOUT;voip_state:IDLE;remain_time:-1;" +
			"imei:000000000000000;imsi:000000000000000;iccid:00000000000000000000;pro:CHINA UNICOM GSM;idle:1;disable_status:0;" +
			"SMS_LOGIN:N;SMB_LOGIN:;CELLINFO:LAC:FFFF,CELL ID:FFFF;CGATT:Y;",
		&RegistrationUpdate{
			RequestID:    "3",
			ID:           "goip-1",
			Password:     "goip",
			Status:       VoIPStatusLogin,
			VoIPStatus:   VoIPStatusLogout,
			VoIPState:    "IDLE",
			VoIPIdleTime: 1,
			IMEI:         "000000000000000",
			IMSI:         "000000000000000",
			ICCID:        "00000000000000000000",
			Carrier:      "CHINA UNICOM GSM",
			CellInfo:     "LAC:FFFF,CELL ID:FFFF",
			CGATT:        true,
			Disabled:     false,
		},
	},
	// region call
	{
		"STATE:1571812793;id:goip-1;password:goip;gsm_remain_state:IDLE",
		&StateUpdate{"1571812793", "goip-1", "goip", "IDLE"},
	},
	{
		"RECORD:1571812801;id:goip-1;password:goip;dir:1;num:1380013800",
		&RecordUpdate{"1571812801", "goip-1", "goip", DirectionIncoming, "1380013800"},
	},
	{
		"RECORD:1571812801;id:goip-1;password:goip;dir:2;num:10086",
		&RecordUpdate{"1571812801", "goip-1", "goip", DirectionOutgoing, "10086"},
	},
	{
		"EXPIRY:1571812871;id:goip-1;password:goip;exp:0",
		&CallTimeUpdate{"1571812871", "goip-1", "goip", nil},
	},
	{
		"HANGUP:1571812814;id:goip-1;password:goip;num:,cause:Normal call clearing",
		&HangupUpdate{"1571812814", "goip-1", "goip", ",cause:Normal call clearing"},
	},
	// TODO: RemainTimeUpdate
	{
		"REMAIN:7568;id:goip-1;password:goip;gsm_remain_time:-1",
		&RemainTimeUpdate{"7568", "goip-1", "goip", nil},
	},
	{
		"GSM:1;get_gsm_num:(null);get_exp_time:0;get_remain_time:0;get_gsm_state:IDLE;imei:000000000000000;out_call_interval:0;module_down:0",
		&GSMUpdate{"1", nil, "IDLE", nil, nil, nil, "000000000000000", false},
	},
	{
		"GSM:1;get_gsm_num:+861380013800;get_exp_time:0;get_remain_time:0;get_gsm_state:IDLE;imei:000000000000000;out_call_interval:0;module_down:0",
		&GSMUpdate{"1", &casePhone, "IDLE", nil, nil, nil, "000000000000000", false},
	},
	// endregion
	// region message
	{
		"RECEIVE:1571812849;id:goip-1;password:goip;srcnum:1380013800;msg:测试",
		&ReceiveSMSUpdate{"1571812849", "goip-1", "goip", "1380013800", "测试"},
	},
	{
		"DELIVER:1571812803;id:goip-1;password:goip;sms_no:32;state:0;num:+861380013800",
		&DeliverSMSUpdate{"1571812803", "goip-1", "goip", "+861380013800", "32", 0},
	},
	// TODO: USSNUpdate
	// endregion
	// region baseband
	{
		"BCCH:34348;id:goip-1;password:goip;bcch:689,638,783,785,640,",
		&BCCHUpdate{"34348", "goip-1", "goip", []string{"689", "638", "783", "785", "640"}},
	},
	{
		"CELLS:34347;id:goip-1;password:goip;lists:689,638,783,785,640,",
		&CellsUpdate{"34347", "goip-1", "goip", []string{"689", "638", "783", "785", "640"}},
	},
	{
		"AT:20;count:1;id:goip-1;password:goip;receive:write:AT+CREG?",
		&ATUpdate{"20", "goip-1", "goip", 1, "write:AT+CREG?"},
	},
	// endregion
}

var badCases = []string{
	"SAMPLE:1;sample",
	"SAMPLE:1;id:nan;",
	"SAMPLE:1;uid:nan;",
}

func TestUnmarshalNotify(t *testing.T) {
	for _, testCase := range testCases {
		expected := reflect.TypeOf(testCase.Expected).Elem()
		t.Run(expected.Name(), func(t *testing.T) {
			message := reflect.New(expected).Interface()
			assert.NoError(t, unmarshalNotifyPacket(testCase.Packet, message))
			assert.Equal(t, testCase.Expected, message)
		})
	}
	for _, packet := range badCases {
		message := new(struct {
			RequestID  string `field:"SAMPLE"`
			ID         int    `field:"id"`
			UnsignedID uint   `field:"uid"`
		})
		err := unmarshalNotifyPacket(packet, message)
		assert.EqualError(t, err, errParseNotifyError.Error())
	}
}

//noinspection SpellCheckingInspection
func TestParseCellInfo(t *testing.T) {
	mapping := map[string]*CellInfo{
		"":                      nil,
		"LAC:FFFF,CELL ID:FFFF": {"FFFF", "FFFF"},
	}
	for input, expected := range mapping {
		info := ParseCellInfo(input)
		assert.Equal(t, expected, info)
	}
}

//noinspection SpellCheckingInspection
func TestOnReplyNotify(t *testing.T) {
	t.Run("Registration", func(t *testing.T) {
		message := &RegistrationUpdate{RequestID: "100"}
		mapping := map[string]error{
			"reg:100;status:0;":  nil,
			"reg:100;status:-1;": errNotFound,
		}
		for expected, err := range mapping {
			reply := onReplyNotify(message, err)
			assert.Equal(t, expected, reply)
		}
	})
	t.Run("General", func(t *testing.T) {
		message := &USSNUpdate{RequestID: "100"}
		mapping := map[string]error{
			"USSN 100 OK":        nil,
			"USSN 100 DISABLE":   ErrDisabledModule,
			"USSN 100 ERROR EOF": io.EOF,
		}
		for expected, err := range mapping {
			reply := onReplyNotify(message, err)
			assert.Equal(t, expected, reply)
		}
	})
	t.Run("NoReply", func(t *testing.T) {
		messages := []interface{}{
			new(ATUpdate),
			new(GSMUpdate),
		}
		for _, message := range messages {
			assert.Empty(t, onReplyNotify(message, nil))
		}
	})
}
