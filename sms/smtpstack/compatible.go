package smtpstack

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func EnableAuthLoginCompatible(be smtp.Backend) smtp.SaslServerFactory {
	return func(conn *smtp.Conn) sasl.Server {
		return sasl.NewLoginServer(func(username, password string) error {
			state := conn.State()
			session, err := be.Login(&state, username, password)
			if err != nil {
				return err
			}
			conn.SetSession(session)
			return nil
		})
	}
}
