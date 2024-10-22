package main

import (
	"bytes"
	"os"

	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/esonhugh/proxyinbrowser/cmd/server/sessionmanager"
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal"
	log "github.com/sirupsen/logrus"
	lef "github.com/t-tomalak/logrus-easy-formatter"
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&lef.Formatter{
		TimestampFormat: "15:04:05",
		LogFormat:       "%time% [id:%session%]> %msg%\n",
	})
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	rch := define.NewChannels()
	go sessionmanager.RunServer(rch, buf)
	// app.RunApp(rch)

	ch := make(chan os.Signal, 1)
	app := terminal.CreateApplication(terminal.ApplicationSpec{
		ConsoleLogBuffer: buf,
		Rch:              rch,
		CloseCh:          ch,
	})
	go app.Run()
	<-ch
	app.Stop()
	os.Exit(0)
}
