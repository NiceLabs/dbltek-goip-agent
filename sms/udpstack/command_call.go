package udpstack

import (
	"time"
)

func (c *Command) GSM() error {
	return c.reply("GSM", holdID, holdPassword)
}

func (c *Command) VoIPStatus() (string, error) {
	return c.get("get_gsm_state")
}

func (c *Command) VoIPTotalTime() (time.Duration, error) {
	return c.getTime("get_exp_time")
}

func (c *Command) SetVoIPTotalTime(t time.Duration) error {
	return c.set("set_exp_time", t)
}

func (c *Command) VoIPRemainTime() (time.Duration, error) {
	return c.getTime("get_remain_time")
}

func (c *Command) ResetVoIPRemainTime(t time.Duration) error {
	return c.execute("reset_remain_time")
}

func (c *Command) OutcallInterval() (time.Duration, error) {
	return c.getTime("get_out_call_interval")
}
func (c *Command) SetOutcallInterval(t time.Duration) error {
	return c.set("set_out_call_interval", t)
}

func (c *Command) EndCall() error {
	return c.execute("svr_drop_call")
}

//CallForward
//  Reason definition
//    0: Unconditional
//    1: Mobile busy
//    2: No reply
//    3: Not reachable
//    4: All call forwarding
//    5: All conditional call forwarding
//  Mode definition
//    0: Disable
//    1: Enable
//    2: Query status
//    3: Registration
//    4: Erasure
func (c *Command) CallForward(reason, mode int, phone string, timeout time.Duration) (status string, err error) {
	id := c.nextID()
	status, err = c.invoke(cmdCF, id, holdPassword, reason, mode, phone, int(timeout.Seconds()))
	if err == nil && mode == 2 {
		err = c.reply(cmdDone, id)
	}
	return
}

//CallWaiting
//  Mode definition
//    0: Disable
//    1: Enable
//    2: Query status
func (c *Command) CallWaiting(mode int) (status string, err error) {
	id := c.nextID()
	status, err = c.invoke(cmdCW, id, holdPassword, mode)
	if err == nil && mode == 2 {
		err = c.reply(cmdDone, id)
	}
	return
}
