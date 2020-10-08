package smtpstack

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
)

var (
	reParseBody = regexp.MustCompile(`^SN:(\w+) Channel:(\d+) Sender:(.+)$`)
)

type session struct {
	Run  func(message *Payload)
	Find func(serial, channel string) string
}

func (s *session) Mail(string, smtp.MailOptions) error { return nil }
func (s *session) Rcpt(string) error                   { return nil }
func (s *session) Reset()                              {}
func (s *session) Logout() error                       { return nil }

func (s *session) Data(r io.Reader) (err error) {
	email, body, err := decodeEmailBody(r)
	if err != nil {
		return
	}
	matched := reParseBody.FindStringSubmatch(body)
	if len(matched) != 4 {
		return
	}
	sender := strings.SplitN(matched[3], ",", 3)
	if len(sender) != 3 {
		return
	}
	go s.Run(&Payload{
		To:           email.To[0].Address,
		SerialNumber: matched[1],
		Channel:      matched[2],
		Phone:        s.Find(matched[1], matched[2]),
		Date:         sender[0],
		Sender:       sender[1],
		Message:      sender[2],
	})
	return
}

func decodeEmailBody(r io.Reader) (email parsemail.Email, body string, err error) {
	if email, err = parsemail.Parse(r); err != nil {
		return
	}
	decoder := base64.NewDecoder(
		base64.StdEncoding,
		strings.NewReader(email.TextBody),
	)
	if decodedBody, err := ioutil.ReadAll(decoder); err == nil {
		body = strings.TrimSpace(string(decodedBody))
	}
	return
}
