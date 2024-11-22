package terminal

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ReOrderArgs(cmd string) []string {
	cmdSplit := strings.Split(cmd, " ")
	var args []string
	for i := 0; i < len(cmdSplit); i++ {
		currentToken := cmdSplit[i]
		if currentToken == "" {
			continue
		}
		if currentToken[0] == '"' {
			for j := i; j < len(cmdSplit); j++ {
				end_token := cmdSplit[j]
				if end_token == "" {
					continue
				}
				if end_token[len(end_token)-1] == '"' {
					currentToken = strings.Join(cmdSplit[i:j+1], " ")
					i = j
					break
				}
			}
			currentToken = strings.TrimPrefix(strings.TrimSuffix(currentToken, "\""), "\"")
		} else if currentToken[0] == '\'' {
			for j := i; j < len(cmdSplit); j++ {
				end_token := cmdSplit[j]
				if end_token == "" {
					continue
				}
				if end_token[len(end_token)-1] == '\'' {
					currentToken = strings.Join(cmdSplit[i:j+1], " ")
					i = j
					break
				}
			}
			currentToken = strings.TrimPrefix(strings.TrimSuffix(currentToken, "'"), "'")
		}
		if currentToken != "" {
			args = append(args, currentToken)
		}
	}
	return args
}

var (
	uid     string
	stopper = make(chan struct{}, 1)
)

func (app Application) ExecuteCommand(cmd string) {
	cmdSplited := ReOrderArgs(cmd)
	// log.Traceln("Cmd: ", "["+strings.Join(cmdSplited, "][")+"]")
	if len(cmdSplited) >= 1 && cmdSplited[0] == "clear" {
		app.Spec.ConsoleLogBuffer.Reset()
		app.LogArea.Clear()
	}
	tcmd.RootCmd.SetArgs(cmdSplited)
	tcmd.RootCmd.SetOut(app.Spec.ConsoleLogBuffer)
	tcmd.RootCmd.CompletionOptions.DisableDefaultCmd = true
	tcmd.RootCmd.DisableSuggestions = true
	err := tcmd.RootCmd.Execute()
	if err != nil {
		log.Error(err)
	}

	/*
		cmdSplited := ReOrderArgs(cmd)
		// log.Traceln("Cmd: ", "["+strings.Join(cmdSplited, "][")+"]")
		var conn *websocket.Conn
		if uid != "" {
			conn = sessionmanager.WebsocketConnMap.Get(uid)
		}
		var taskID = uuid.New().String()
		log := log.WithField("session", uid)
		// System Command
		if len(cmdSplited) < 0 {
			return
		}
		{
			switch cmdSplited[0] {
			case "help":
				log.Infoln("\nCommands: \n" +
					"  uid: Show current session id\n" +
					"  list: List all session id\n" +
					"  use [uid]: Use session id, if not specified, use first one\n" +
					"  relay: Start HTTP/S Proxy to relay any request to remote browser\n" +
					"  stop: Stop HTTP/S Proxy\n" +
					"  clear: clean screen" +
					"  exit/quit: Exit")
				return
			case "uid":
				log.Infoln("current session id: ", uid)
				return
			case "list":
				list := sessionmanager.WebsocketConnMap.List()
				log.Infoln("\n======LIST======\n" + strings.Join(list, "\n"))
				return
			case "use":
				if len(cmdSplited) < 2 {
					list := sessionmanager.WebsocketConnMap.List()
					if len(list) == 0 {
						log.Errorln("No session id available")
						return
					} else {
						uid = list[0]
						log.Infoln("Use default uid: ", uid)
					}
					return
				}
				uid = cmdSplited[1]
				log.Infoln("Use ", uid)
				return
			case "clear":
				app.Spec.ConsoleLogBuffer.Reset()
				app.LogArea.Clear()
				return
			case "quit", "exit":
				app.Spec.CloseCh <- syscall.SIGTERM
				os.Exit(0)
			}
		}

		{
			// Task command
			if conn == nil {
				log.Errorln("Connection not found or uid error")
				return
			}
			switch cmdSplited[0] {
			case "relay":
				go http_proxy.CreateHttpProxyServer(conn, "9001", app.Spec.Rch, stopper)
			case "stop":
				stopper <- struct{}{}
				log.Infof("SendTo Stop signal to stop relay\n")
				stopper = make(chan struct{}, 1)
			}
			log.Infoln("Beacon command sent. Task: ", taskID)
		}
	*/
}
