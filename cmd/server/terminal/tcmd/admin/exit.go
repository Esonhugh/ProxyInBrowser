package admin

import (
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	tcmd.RootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:     "exit",
	Aliases: []string{"quit", "q"},
	Short:   "Exit Application",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}
