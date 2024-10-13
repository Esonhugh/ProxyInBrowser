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
		LogFormat:       "%time%[%lvl%]> %msg%\n",
	})
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	rch := define.NewChannels()
	go sessionmanager.RunServer(rch, buf)
	// app.RunApp(rch)

	ch := make(chan os.Signal, 1)
	go terminal.RunApplication(terminal.ApplicationSpec{
		ConsoleLogBuffer: buf,
		Rch:              rch,
		CloseCh:          ch,
	})
	<-ch
	os.Exit(0)
}
