package main

import (
	"net"
	"strings"

	"github.com/rs/zerolog/log"
)

type PacketConnProxy struct {
	net.PacketConn
}

func (c *PacketConnProxy) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = c.PacketConn.ReadFrom(p)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return n, addr, err
	}
	message := string(p[:n])
	event := log.Info()
	if strings.HasPrefix(message, "req:") {
		event = log.Debug()
	}
	if err != nil {
		event.Err(err)
	} else {
		event.Str("direction", "incoming").Str("remote", addr.String()).Msg(string(p[:n]))
	}
	return
}

func (c *PacketConnProxy) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = c.PacketConn.WriteTo(p, addr)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return n, err
	}
	message := string(p[:n])
	event := log.Info()
	if strings.HasPrefix(message, "reg:") {
		event = log.Debug()
	}
	if err != nil {
		event.Err(err)
	} else {
		event.Str("direction", "outgoing").Str("remote", addr.String()).Msg(string(p[:n]))
	}
	return
}
