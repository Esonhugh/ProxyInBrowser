package admin

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/sessionmanager"
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	tcmd.RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all available sessions",
	Run: func(cmd *cobra.Command, args []string) {
		list := sessionmanager.WebsocketConnMap.List()
		if len(list) == 0 {
			tcmd.Opt.Log.Infof("no sessions found")
		}
		tcmd.Opt.Log.Infoln("\n======LIST======\n" + strings.Join(list, "\n") + "\n")
	},
}
