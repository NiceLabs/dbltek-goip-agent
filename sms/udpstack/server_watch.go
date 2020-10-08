package udpstack

import (
	"context"
	"net"
	"regexp"
	"strings"
	"time"
)

var (
	reNotifyName = regexp.MustCompile(`^(\w{1,10}):`)
	reRequestID  = regexp.MustCompile(`^\w+ (\d+) `)
)

func (s *Server) Watch(ctx context.Context) {
	packet := make([]byte, 0x8000) // 32 KiB
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_ = s.conn.SetDeadline(time.Now().Add(time.Second * 30))
			n, addr, err := s.conn.ReadFrom(packet)
			if err != nil {
				continue
			}
			go s.onHandler(addr, strings.TrimSpace(string(packet[0:n])))
		}
	}
}

func (s *Server) onHandler(addr net.Addr, packet string) {
	if reply := s.onNotifyHandler(addr, packet); reply != "" {
		_, _ = s.conn.WriteTo([]byte(reply), addr)
		return
	}
	var id string
	if matched := reRequestID.FindStringSubmatch(packet); len(matched) == 2 {
		id = matched[1]
	}
	if id == "" {
		return
	} else if err := s.onInvokeHandler(addr, id, packet); err == nil {
		return
	} else if err := s.onBulkSMSHandler(addr, id, packet); err == nil {
		return
	}
}

func (s *Server) onNotifyHandler(addr net.Addr, packet string) (reply string) {
	if s.OnNotify == nil {
		return
	}
	//noinspection SpellCheckingInspection
	mapping := map[string]interface{}{
		"req":     new(RegistrationUpdate), // Keep-Alive: req
		"RECEIVE": new(ReceiveSMSUpdate),   //        SMS: RECEIVE
		"DELIVER": new(DeliverSMSUpdate),   //        SMS: DELIVER
		"USSN":    new(USSNUpdate),         //        SMS: USSN
		"STATE":   new(StateUpdate),        //       VoIP: STATE
		"RECORD":  new(RecordUpdate),       //       VoIP: RECORD
		"HANGUP":  new(HangupUpdate),       //       VoIP: HANGUP
		"REMAIN":  new(RemainTimeUpdate),   //       VoIP: REMAIN
		"EXPIRY":  new(CallTimeUpdate),     //       VoIP: EXPIRY
		"GSM":     new(GSMUpdate),          //       VoIP: GSM
		"BCCH":    new(BCCHUpdate),         //   Baseband: BCCH
		"CELLS":   new(CellsUpdate),        //   Baseband: CELLS
		"AT":      new(ATUpdate),           //   Baseband: AT
	}
	var name string
	if matched := reNotifyName.FindStringSubmatch(packet); len(matched) == 2 {
		name = matched[1]
	}
	if message, ok := mapping[name]; !ok {
		return
	} else if err := unmarshalNotifyPacket(packet, message); err != nil {
		return onReplyNotify(message, err)
	} else {
		return onReplyNotify(message, s.OnNotify.OnNotify(addr, message))
	}
}

func (s *Server) onInvokeHandler(addr net.Addr, id, packet string) error {
	if handle, ok := s.pending[id]; ok {
		handle(packet)
		return nil
	}
	return errNotFound
}

func (s *Server) onBulkSMSHandler(addr net.Addr, id, packet string) error {
	if s.OnBulkSMS == nil || !s.OnBulkSMS.CanSent(id) {
		return errNotFound
	}
	var (
		handler  = s.OnBulkSMS
		commands = strings.Split(packet, cmdJoinSymbol)
		name     = commands[0]
		onDone   = func() { _ = s.reply(addr, []string{cmdDone, id}) }
	)
	switch name {
	case cmdPassword:
		password, err := handler.OnPassword(id)
		if err != nil {
			onDone()
		} else {
			_ = s.reply(addr, []string{cmdPassword, commands[1], password})
		}
	case cmdSend, cmdOK, cmdError:
		if len(commands) == 4 {
			var err error
			if name == cmdOK {
				index := commands[2]
				deliverID := commands[3]
				err = handler.OnSendSent(id, index, deliverID)
			} else if name == cmdError {
				index := commands[2]
				message := commands[3]
				err = handler.OnSendFailed(id, index, message)
			}
			if err != nil {
				onDone()
				return errNotFound
			}
		}
		index, phone, err := handler.OnSend(id)
		if err != nil {
			onDone()
		} else {
			_ = s.reply(addr, []string{cmdSend, id, index, phone})
		}
	case cmdWait:
		index := commands[2]
		if err := handler.OnSendWait(id, index); err != nil {
			onDone()
			return errNotFound
		}
	case cmdDone:
		_ = handler.OnEnd(id)
	}
	return nil
}
