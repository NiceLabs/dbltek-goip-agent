package udpstack

import (
	"net"
)

//noinspection SpellCheckingInspection
type Server struct {
	conn      net.PacketConn
	pending   map[string]func(packet string)
	OnNotify  OnNotifyHandler
	OnBulkSMS OnBulkSMSHandler
	OnUSSD    OnUSSDHandler
}

func NewServer(conn net.PacketConn) *Server {
	return &Server{
		conn:    conn,
		pending: make(map[string]func(string)),
	}
}

func (s *Server) Command(addr net.Addr, slotId uint, password string) *Command {
	return &Command{
		server:   s,
		addr:     addr,
		slotID:   slotId,
		password: password,
	}
}
