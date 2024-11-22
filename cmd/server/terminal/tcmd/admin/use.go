package admin

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/sessionmanager"
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	tcmd.RootCmd.AddCommand(UseCmd)
}

var UseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a specific session or change into it",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			list := sessionmanager.WebsocketConnMap.List()
			if len(list) == 0 {
				log.Errorln("No session id available")
				return
			} else {
				uid := list[0]
				tcmd.Opt.SessionId = uid
			}
		} else {
			tcmd.Opt.SessionId = args[0]
		}
		log.Infoln("Change to uid: ", tcmd.Opt.SessionId)
	},
}
