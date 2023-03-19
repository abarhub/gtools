package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/rm"
)

var (
//longFormat bool
//recursive  bool
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
		param := rm.RmParameters{path, longFormat, recursive}
		commandError = cmdRm(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func cmdRm(param rm.RmParameters) error {
	return rm.RmCommand(param)
}

func ConfigureRmCommandLine(rootCmd *cobra.Command) {

	//rmCmd.Flags().BoolVarP(&longFormat, "", "l", false,
	//	"Long listing format")
	//rmCmd.Flags().BoolVarP(&recursive, "recursive", "r", false,
	//	"List subdirectory recursively")
	rootCmd.AddCommand(rmCmd)

}
