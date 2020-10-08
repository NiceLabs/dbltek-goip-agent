package udpstack

import (
	"net"
	"strings"
)

type Command struct {
	server   *Server
	addr     net.Addr
	slotID   uint
	password string
}

func (c *Command) Reply(command string) error {
	id := c.nextID()
	command = strings.Replace(command, "{id}", id, -1)
	command = strings.Replace(command, "{password}", c.password, -1)
	return c.reply(command)
}

func (c *Command) RebootModule() error {
	return c.execute("svr_reboot_module")
}

func (c *Command) RebootDevice() error {
	return c.execute("svr_reboot_dev")
}

func (c *Command) Module(enabled bool) error {
	state := 2 // disabled
	if enabled {
		state = 1
	}
	return c.set("module_ctl_i", state)
}
