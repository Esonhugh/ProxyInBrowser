package ops

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/http_proxy"
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
)

var port string

func init() {
	relayCmd.PersistentFlags().StringVarP(&port, "port", "p", "9001", "Port to listen on")
	tcmd.RootCmd.AddCommand(relayCmd)
}

var relayCmd = &cobra.Command{
	Use:     "relay",
	Short:   "start a relay session",
	PreRunE: PreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		go http_proxy.Serve(tcmd.Opt.Session, port)
	},
}

var stopRelayCmd = &cobra.Command{
	Use:     "stop",
	Short:   "stop a relay session",
	PreRunE: PreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		http_proxy.Stop()
		tcmd.Opt.Log.Infof("SendTo Stop signal to stop relay\n")
	},
}
