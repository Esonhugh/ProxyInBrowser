package admin

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
)

func init() {
	tcmd.RootCmd.AddCommand(UidCmd)
}

var UidCmd = &cobra.Command{
	Use:   "uid",
	Short: "Lookup current uid",
	Run: func(cmd *cobra.Command, args []string) {
		if tcmd.Opt.SessionId == "" || tcmd.Opt.Session == nil {
			tcmd.Opt.Log.Infof("session id useless")
		} else {
			tcmd.Opt.Log.Infof("Current session id is %v ", tcmd.Opt.SessionId)
		}
	},
}
