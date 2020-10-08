package udpstack

func (c *Command) GSMNumber() (string, error) {
	return c.get("get_gsm_number")
}

func (c *Command) SetGSMNumber(code string) error {
	return c.set("set_gsm_number", code)
}

func (c *Command) IMEI() (string, error) {
	return c.get("get_imei")
}

func (c *Command) Cells() error {
	return c.execute("get_cells_list")
}

func (c *Command) CurrentCell() (string, error) {
	return c.get("CURCELL")
}

func (c *Command) SetBaseCell(cell string) error {
	return c.set("set_base_cell", cell)
}

func (c *Command) SetIMEI(code string) error {
	return c.set("set_imei", code)
}

func (c *Command) ATStart() error {
	return c.execute("ATSTART")
}

func (c *Command) ATStop() error {
	return c.execute("ATSTOP")
}

func (c *Command) ATAlive() error {
	return c.execute("ATALIVE")
}

func (c *Command) ATCommand(command string) error {
	_, err := c.invoke("ATCMD", holdID, command, holdPassword)
	return err
}
