package smtpstack

// L1 DeviceSN
// L2 Channel
// L3 Phone number
type PhoneMapping map[string]map[string]string

type Configuration struct {
	SMTPUsername string       `json:"smtp_username"`
	SMTPPassword string       `json:"smtp_password"`
	Hook         string       `json:"hook"`
	Phones       PhoneMapping `json:"phones"`
}
