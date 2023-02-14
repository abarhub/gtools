package utils

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Command int64

const (
	Copy Command = iota
	Du
)

type Option int64

const (
	Include Option = iota
	Exclude
)

type Parameters struct {
	Command   Command
	Arguments []string
	Options   map[Option]string
}

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

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		commandError = cmdCopy(args[0], args[1])
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
		commandError = cmdDu(path, humanReadable, thresholdStr, maxDepth)
		return nil
	},
}

func Run(args []string) error {

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

func cmdCopy(src string, dest string) error {
	return CopyDir(src, dest)
}

func cmdDu(path string, humanReadable bool, thresholdStr string, maxDepth int) error {
	return DiskUsage(path, humanReadable, thresholdStr, maxDepth)
}
