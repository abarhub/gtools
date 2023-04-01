package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/ls"
)

var (
	longFormat         bool
	recursive          bool
	excludePathLs      []string
	includePathLs      []string
	displayDirectoryLs bool
	displayAllLs       bool
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
		param := ls.LsParameters{path, longFormat, recursive,
			excludePathLs, includePathLs, displayDirectoryLs, displayAllLs}
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
	lsCmd.Flags().StringSliceVarP(&excludePathLs, "exclude", "e", []string{}, "Path to exclude")
	lsCmd.Flags().StringSliceVarP(&includePathLs, "include", "i", []string{}, "Path to include")
	lsCmd.Flags().BoolVarP(&displayDirectoryLs, "displayDirectory", "d", true,
		"Display directory")
	lsCmd.Flags().BoolVarP(&displayAllLs, "displayAll", "a", false,
		"Display dot files")
	rootCmd.AddCommand(lsCmd)

}
