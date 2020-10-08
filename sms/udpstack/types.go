package udpstack

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrDisabledModule   = errors.New("udpstack: Disabled Module")
	ErrUpdateFailed     = errors.New("udpstack: Update failed")
	ErrUnboundOnBulkSMS = errors.New("udpstack: OnBulkSMS unbound instance")
	errParseNotifyError = errors.New("udpstack: Parse notify packet error")
	errNotFound         = errors.New("udpstack: Not found")
)

type (
	OnNotifyHandler interface {
		OnNotify(addr net.Addr, update interface{}) error
	}
	OnUSSDHandler interface {
		OnUSSD(slotID uint, command, message string, err error)
	}
	OnBulkSMSHandler interface {
		CanSent(id string) (ok bool)
		OnStart(slotID uint, id string, phones []string, message string) (err error)
		OnPassword(id string) (password string, err error)
		OnSend(id string) (index, phone string, err error)
		OnSendFailed(id, index, message string) (err error)
		OnSendWait(id, index string) (err error)
		OnSendSent(id, index, deliverID string) (err error)
		OnEnd(id string) (err error)
	}
)

// region direction

type Direction int

const (
	DirectionIncoming Direction = iota + 1
	DirectionOutgoing
)

func (d Direction) String() string {
	switch d {
	case DirectionIncoming:
		return "incoming"
	case DirectionOutgoing:
		return "outgoing"
	}
	return ""
}

// endregion

// region VoIP status

type VoIPStatus string

const (
	VoIPStatusLogin  VoIPStatus = "LOGIN"
	VoIPStatusLogout            = "LOGOUT"
	VoIPStatusUp                = "UP"
	VoIPStatusDown              = "DOWN"
)

// endregion

// region error

type Error string

//noinspection SpellCheckingInspection
func (p Error) Error() string {
	splitted := strings.Split(string(p), cmdJoinSymbol)
	if len(splitted) < 3 {
		return ""
	}
	return splitted[len(splitted)-1]
}

// endregion
