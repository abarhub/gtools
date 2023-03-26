package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/rm"
)

var (
	confirmationRm bool
	recursiveRm    bool
	verboseRm      bool
	excludePathRm  []string
	includePathRm  []string
	dryRunRm       bool
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
		param := rm.RmParameters{path, confirmationRm, recursiveRm,
			verboseRm, excludePathRm, includePathRm, dryRun}
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
	rmCmd.Flags().StringSliceVarP(&excludePath, "exclude", "e", []string{}, "Path to exclude")
	rmCmd.Flags().StringSliceVarP(&includePath, "include", "i", []string{}, "Path to include")
	rmCmd.Flags().BoolVarP(&dryRun, "dryRun", "d", false, "Don't remove")
	rootCmd.AddCommand(rmCmd)

}
