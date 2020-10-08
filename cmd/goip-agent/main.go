package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/NiceLabs/dbltek-goip-agent/sms/api"
	"github.com/NiceLabs/dbltek-goip-agent/sms/business"
	"github.com/NiceLabs/dbltek-goip-agent/sms/udpstack"
)

var (
	BuildInfo = "Development edition"
)

//noinspection GoUnhandledErrorResult
func main() {
	fmt.Println("Build Info: ", BuildInfo)

	engine := openDatabase()
	defer engine.Close()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	setLogger()

	server := listenServer(":44444")
	business.BindDatabase(engine, server)
	go server.Watch(context.Background())

	handler := api.NewRPCHandler(engine, server)
	log.Print(http.ListenAndServe(":10000", handler))
}

func openDatabase() *gorm.DB {
	// see https://github.com/mattn/go-sqlite3#faq
	engine, err := gorm.Open("sqlite3", "agent.db?cache=shared&mode=rwc")
	if err != nil {
		panic(err)
	}
	engine.DB().SetMaxOpenConns(1)
	return engine
}

func setLogger() {
	fp, err := os.OpenFile("agent.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(fp)
}

func listenServer(address string) *udpstack.Server {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		panic(err)
	}
	conn = &PacketConnProxy{
		PacketConn: conn,
	}
	return udpstack.NewServer(conn)
}
