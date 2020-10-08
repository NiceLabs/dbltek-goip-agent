package api

import (
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/jinzhu/gorm"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

func NewRPCHandler(engine *gorm.DB, serve *udpstack.Server) http.Handler {
	handler := rpc.NewServer()
	handler.RegisterCodec(json2.NewCodec(), "application/json")
	_ = handler.RegisterService(&BaseService{engine, serve}, "Base")
	return handler
}
