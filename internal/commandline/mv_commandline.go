package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
	"gtools/internal/command/mv"
)

var (
	copyAndDelete bool
	verboseMv     bool
)

var mvCmd = &cobra.Command{
	Use:   "mv [flags] src dest",
	Short: "move files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		pathSrc := ""
		pathDest := ""
		if len(args) == 2 {
			pathSrc = args[0]
			pathDest = args[1]
		} else {
			return fmt.Errorf("mv need 2 arguments")
		}
		param := mv.MvParameters{pathSrc, pathDest, copyAndDelete, verboseMv}
		commandError = cmdMv(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func cmdMv(param mv.MvParameters) error {
	return mv.MvCommand(param)
}

func ConfigureMvCommandLine(rootCmd *cobra.Command) {

	mvCmd.Flags().BoolVarP(&copyAndDelete, "copyAndDelete", "c", false,
		"Copy to destination and delete source (only for file move)")
	mvCmd.Flags().BoolVarP(&verboseMv, "verbose", "v", false,
		"Verbose")
	rootCmd.AddCommand(mvCmd)

}
