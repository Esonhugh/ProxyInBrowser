package main

import (
	"bytes"
	"os"

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
	go sessionmanager.RunServer(buf)
	// app.RunApp(rch)

	ch := make(chan os.Signal, 1)
	app := terminal.CreateApplication(terminal.ApplicationSpec{
		ConsoleLogBuffer: buf,
		CloseCh:          ch,
	})
	go app.Run()
	<-ch
	app.Stop()
	os.Exit(0)
}
