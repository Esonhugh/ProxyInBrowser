package admin

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
)

func init() {
	tcmd.RootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear screen",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
