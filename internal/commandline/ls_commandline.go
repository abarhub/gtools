package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/ls"
)

var (
	longFormat bool
	recursive  bool
)

var lsCmd = &cobra.Command{
	Use:   "ls [flags] [dir]",
	Short: "list files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		param := ls.LsParameters{path, longFormat, recursive}
		commandError = cmdLs(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func cmdLs(param ls.LsParameters) error {
	return ls.LsCommand(param)
}

func ConfigureLsCommandLine(rootCmd *cobra.Command) {

	lsCmd.Flags().BoolVarP(&longFormat, "", "l", false,
		"Long listing format")
	lsCmd.Flags().BoolVarP(&recursive, "recursive", "r", false,
		"List subdirectory recursively")
	rootCmd.AddCommand(lsCmd)

}
