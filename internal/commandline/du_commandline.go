package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/du"
)

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
		commandError = cmdDu(du.DuParameters{path, humanReadable, thresholdStr,
			maxDepth, []string{}, []string{}})
		return nil
	},
}

func cmdDu(param du.DuParameters) error {
	return du.DiskUsage(param)
}

func ConfigureDuCommandLine(rootCmd *cobra.Command) {

	duCmd.Flags().BoolVarP(&humanReadable, "humanReadable", "r", false, "\"Human-readable\" output.  Use unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte.")
	duCmd.Flags().StringVarP(&thresholdStr, "threshold", "t", "", "threshold of the size, any folders' size larger than the threshold will be print. for example, '1G', '10M', '100K', '1024'")
	duCmd.Flags().IntVarP(&maxDepth, "maxDepth", "d", 0, "list its subdirectories and their sizes to any desired level of depth (i.e., to any level of subdirectories) in a directory tree.")

	rootCmd.AddCommand(duCmd)
}
