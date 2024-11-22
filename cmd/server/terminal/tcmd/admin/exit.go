package admin

import (
	"github.com/spf13/cobra"
	"os"
)

var exitCmd = &cobra.Command{
	Use:     "exit",
	Aliases: []string{"quit", "q"},
	Short:   "Exit Application",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}
