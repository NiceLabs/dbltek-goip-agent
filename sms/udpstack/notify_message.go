package udpstack

//noinspection SpellCheckingInspection
type (
	ReceiveSMSUpdate struct {
		RequestID string `field:"RECEIVE"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Phone     string `field:"srcnum"`
		Message   string `field:"msg"`
	}
	DeliverSMSUpdate struct {
		RequestID string `field:"DELIVER"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Phone     string `field:"num"`
		SMSNo     string `field:"sms_no"`
		State     int    `field:"state"`
	}
	USSNUpdate struct {
		RequestID string `field:"USSN"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Message   string `field:"msg"`
	}
)
