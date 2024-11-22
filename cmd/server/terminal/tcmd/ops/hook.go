package ops

import (
	"errors"
	"github.com/esonhugh/proxyinbrowser/cmd/server/terminal/tcmd"
	"github.com/spf13/cobra"
)

func PreRunE(cmd *cobra.Command, args []string) error {
	if tcmd.Opt.Session == nil {
		return errors.New("No session specified. ")
	}
	return nil
}

func PostRunE(cmd *cobra.Command, args []string) error {
	tcmd.Opt.Log.Infof("Command %v Sent", tcmd.Opt.TaskId)
	return nil
}
