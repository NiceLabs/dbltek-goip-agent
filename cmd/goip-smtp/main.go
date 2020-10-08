package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/NiceLabs/dbltek-goip-agent/sms/smtpstack"
)

var configure = new(smtpstack.Configuration)

func init() {
	configFile, err := ioutil.ReadFile(path.Join(getProgramDirectory(), "configure.json"))
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configFile, configure)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	be := &smtpstack.Backend{
		Configuration: configure,
		RunHook:       runProgram,
	}
	s := smtp.NewServer(be)
	s.Addr = ":25"
	s.AllowInsecureAuth = true
	s.EnableAuth(sasl.Login, smtpstack.EnableAuthLoginCompatible(be))

	log.Println("Starting server at", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func runProgram(message *smtpstack.Payload) {
	log.Print(message.String(), " Started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, path.Join(getProgramDirectory(), configure.Hook))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		//noinspection GoUnhandledErrorResult
		defer stdin.Close()
		encoded, err := json.Marshal(message)
		if err != nil {
			return
		}
		_, _ = stdin.Write(encoded)
	}()
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(message, err, "\n", string(output))
	} else {
		log.Print(message.String(), " End")
	}
}

func getProgramDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
