package udpstack

import "strconv"

func (c *Command) SendMessage(id string, phones []string, message string) (err error) {
	handler := c.server.OnBulkSMS
	if handler == nil {
		return ErrUnboundOnBulkSMS
	}
	if _, err = strconv.ParseInt(id, 10, 31); err != nil {
		return
	}
	if err = handler.OnStart(c.slotID, id, uniqueString(phones), message); err == nil {
		err = c.reply(cmdMessage, id, len(message), message)
	}
	return
}

func (c *Command) USSD(command string) (result string, err error) {
	result, err = c.invoke(cmdUSSD, holdID, holdPassword, command)
	c.server.OnUSSD.OnUSSD(c.slotID, command, result, err)
	return
}

func (c *Command) USSDExit() (err error) {
	return c.execute(cmdUSSDExit)
}

func uniqueString(input []string) (output []string) {
	keys := make(map[string]bool)
	for _, phone := range input {
		if _, value := keys[phone]; !value {
			keys[phone] = true
			output = append(output, phone)
		}
	}
	return
}
