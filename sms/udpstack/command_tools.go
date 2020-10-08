package udpstack

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//noinspection SpellCheckingInspection
const (
	cmdCW        = "CW"
	cmdCF        = "CF"
	cmdMessage   = "MSG"
	cmdPassword  = "PASSWORD"
	cmdSend      = "SEND"
	cmdWait      = "WAIT"
	cmdOK        = "OK"
	cmdError     = "ERROR"
	cmdDone      = "DONE"
	cmdUSSD      = "USSD"
	cmdUSSDExit  = "USSDEXIT"
	holdID       = "\x00ID\x00"
	holdPassword = "\x00PASSWORD\x00"
)

func (c *Command) invoke(input ...interface{}) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return c.server.invoke(ctx, c.addr, c.invokeToCommand(input))
}

func (c *Command) reply(input ...interface{}) error {
	return c.server.reply(c.addr, c.invokeToCommand(input))
}

func (c *Command) invokeToCommand(inputs []interface{}) (outputs []string) {
	id := c.nextID()
	for _, input := range inputs {
		output := ""
		switch v := input.(type) {
		case nil:
			continue
		case string:
			output = v
			switch v {
			case holdID:
				output = id
			case holdPassword:
				output = c.password
			}
		case int:
			output = strconv.Itoa(v)
		case time.Duration:
			output = strconv.Itoa(int(v.Minutes()))
		}
		outputs = append(outputs, output)
	}
	return outputs
}

func (c *Command) get(name string) (string, error) {
	return c.invoke(name, holdID, holdPassword)
}

func (c *Command) getTime(name string) (time.Duration, error) {
	if value, err := c.get(name); err != nil {
		return 0, err
	} else if n, err := strconv.ParseInt(value, 10, 64); err != nil {
		return 0, err
	} else {
		duration := time.Duration(n) * time.Minute
		return duration, nil
	}
}

func (c *Command) set(name string, input interface{}) error {
	result, err := c.invoke(name, holdID, input, holdPassword)
	if strings.ToUpper(result) == cmdOK {
		return nil
	}
	if err != nil {
		return err
	}
	return ErrUpdateFailed
}

func (c *Command) execute(name string) error {
	return c.set(name, nil)
}

func (c *Command) nextID() string {
	return strconv.Itoa(int(rand.Int31()))
}
