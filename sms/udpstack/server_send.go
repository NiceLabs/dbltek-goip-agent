package udpstack

import (
	"context"
	"net"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	reInvokeError = regexp.MustCompile(`^\w+ERROR `)
)

func (s *Server) reply(addr net.Addr, input []string) error {
	command := strings.Join(input, cmdJoinSymbol)
	_, err := s.conn.WriteTo([]byte(command), addr)
	return err
}

//noinspection SpellCheckingInspection
func (s *Server) invoke(ctx context.Context, addr net.Addr, input []string) (value string, err error) {
	if err = s.reply(addr, input); err != nil {
		return
	}
	requestID := input[1]
	reply := make(chan string, 1)
	s.pending[requestID] = func(packet string) {
		delete(s.pending, requestID)
		reply <- packet
	}
	select {
	case packet := <-reply:
		if reInvokeError.MatchString(packet) {
			err = Error(packet)
		} else {
			splitted := strings.Split(packet, cmdJoinSymbol)
			if len(splitted) > 3 {
				value = splitted[len(splitted)-1]
			}
		}
	case <-ctx.Done():
		delete(s.pending, requestID)
		close(reply)
		err = ctx.Err()
		log.Error().
			Str("action", "invoke").
			Strs("input", input).
			Err(err)
	}
	return
}
