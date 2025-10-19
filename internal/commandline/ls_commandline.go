package commandline

import (
	"gtools/internal/command/ls"
	"gtools/internal/utils"

	"github.com/spf13/cobra"
)

var (
	longFormat         bool
	recursive          bool
	excludePathLs      []string
	includePathLs      []string
	displayDirectoryLs bool
	displayAllLs       bool
	datesLs            []string
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
		date2, err := utils.ParseDates(datesLs)
		if err != nil {
			return err
		}
		param := ls.LsParameters{Path: path, LongFormat: longFormat, Recurvive: recursive,
			ExcludePath: excludePathLs, IncludePath: includePathLs, DisplayDirectory: displayDirectoryLs,
			DisplayAll: displayAllLs, Dates: date2}
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
	lsCmd.Flags().StringSliceVarP(&datesLs, "dates", "t", []string{},
		"Date of files to show in format YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS with optional prefix >, <, >=, <=, =")
	rootCmd.AddCommand(lsCmd)

}
