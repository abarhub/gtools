package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/time"
)

var timeCmd = &cobra.Command{
	Use:   "time [flags] cmd [[[param1] param2] ...]",
	Short: "time execution of command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := time.TimeParameters{args[0], args[1:]}
		commandError = time.TimeCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureTimeCommandLine(rootCmd *cobra.Command) {
	rootCmd.AddCommand(timeCmd)
}
