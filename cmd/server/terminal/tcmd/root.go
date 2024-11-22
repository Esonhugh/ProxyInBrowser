package tcmd

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/esonhugh/proxyinbrowser/cmd/server/sessionmanager"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opt struct {
	SessionId string
	Session   *define.WebsocketClient
	TaskId    string

	Log *log.Entry
}

var RootCmd = &cobra.Command{
	Use:   "",
	Short: "Begin of the command execute in console",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Opt.Log = log.WithField("session", Opt.SessionId)
		if Opt.SessionId != "" {
			Opt.Session = sessionmanager.WebsocketConnMap.Get(Opt.SessionId)
		}
		Opt.TaskId = uuid.New().String()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
