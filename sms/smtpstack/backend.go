package smtpstack

import (
	"github.com/emersion/go-smtp"
)

type Backend struct {
	*Configuration
	RunHook func(message *Payload)
}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if bkd.SMTPUsername != username && bkd.SMTPPassword != password {
		return nil, smtp.ErrAuthUnsupported
	}
	sess := &session{
		Run: bkd.RunHook,
		Find: func(serial, channel string) (phone string) {
			if bkd.Phones == nil {
				return
			}
			device, ok := bkd.Phones[serial]
			if !ok {
				return
			}
			phone, ok = device[channel]
			return
		},
	}
	return sess, nil
}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}
