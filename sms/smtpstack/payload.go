package smtpstack

import "fmt"

type Payload struct {
	To           string `json:"to"`
	SerialNumber string `json:"sn"`
	Channel      string `json:"channel"`
	Phone        string `json:"phone,omitempty"`
	Date         string `json:"date"`
	Sender       string `json:"sender"`
	Message      string `json:"message"`
}

func (m *Payload) String() string {
	return fmt.Sprintf("%s#%s", m.SerialNumber, m.Channel)
}
