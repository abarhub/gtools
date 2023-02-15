package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
	copy2 "gtools/internal/command/copy"
	"gtools/internal/command/du"
)

var commandError error = nil

var rootCmd = &cobra.Command{
	Use:   "gtools",
	Short: "gtools - a simple CLI tools",
	Long: `gtools is a super simple CLI tools
   
simple in CLI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("no command")
	},
}

var (
	excludePath   []string
	includePath   []string
	createDestDir bool
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := copy2.CopyParameters{args[0], args[1],
			excludePath, includePath, createDestDir}
		commandError = cmdCopy(param)
		return nil
	},
}

var (
	humanReadable = false
	thresholdStr  = ""
	maxDepth      = 0
)

var duCmd = &cobra.Command{
	Use:   "du",
	Short: "disk usage",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path = ""
		if len(args) > 0 {
			path = args[0]
		}
		commandError = cmdDu(du.DuParameters{path, humanReadable, thresholdStr, maxDepth})
		return nil
	},
}

func Run(args []string) error {

	copyCmd.Flags().StringArrayVarP(&excludePath, "exclude", "e", []string{}, "Path to exclude")
	copyCmd.Flags().StringArrayVarP(&includePath, "include", "i", []string{}, "Path to include")
	copyCmd.Flags().BoolVarP(&createDestDir, "createDestDir", "c", false, "Create destination directory if not exists")
	rootCmd.AddCommand(copyCmd)

	duCmd.Flags().BoolVarP(&humanReadable, "humanReadable", "r", false, "\"Human-readable\" output.  Use unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte.")
	duCmd.Flags().StringVarP(&thresholdStr, "threshold", "t", "", "threshold of the size, any folders' size larger than the threshold will be print. for example, '1G', '10M', '100K', '1024'")
	duCmd.Flags().IntVarP(&maxDepth, "maxDepth", "d", 0, "list its subdirectories and their sizes to any desired level of depth (i.e., to any level of subdirectories) in a directory tree.")

	rootCmd.AddCommand(duCmd)

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("Whoops. There was an error while executing your CLI '%s'", err)
	}

	return commandError
}

func cmdCopy(param copy2.CopyParameters) error {
	return copy2.CopyDir(param)
}

func cmdDu(param du.DuParameters) error {
	return du.DiskUsage(param)
}
