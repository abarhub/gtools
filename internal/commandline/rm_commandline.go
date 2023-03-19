package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/rm"
)

var (
	confirmationRm bool
	recursiveRm    bool
	verboseRm      bool
)

var rmCmd = &cobra.Command{
	Use:   "rm [flags] [dir]",
	Short: "remove files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		param := rm.RmParameters{path, confirmationRm, recursiveRm, verboseRm}
		commandError = cmdRm(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func cmdRm(param rm.RmParameters) error {
	return rm.RmCommand(param)
}

func ConfigureRmCommandLine(rootCmd *cobra.Command) {

	rmCmd.Flags().BoolVarP(&confirmationRm, "confirmation", "c", true,
		"Confirmation")
	rmCmd.Flags().BoolVarP(&recursiveRm, "recursive", "r", false,
		"Remove subdirectory recursively")
	rmCmd.Flags().BoolVarP(&verboseRm, "verbose", "v", false,
		"Verbose")
	rootCmd.AddCommand(rmCmd)

}
