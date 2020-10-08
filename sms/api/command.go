package api

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/NiceLabs/dbltek-goip-agent/sms/business"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

type BaseService struct {
	db    *gorm.DB
	serve *udpstack.Server
}

type TestArgs struct {
	DeviceID string
	Command  string
}

type EmptyReply struct {
	Error string
}

func (b *BaseService) TestSend(r *http.Request, args *TestArgs, reply *EmptyReply) error {
	cmd, err := business.BindCommand(b.db, b.serve, args.DeviceID)
	if err != nil {
		return err
	}
	return cmd.Reply(args.Command)
}
