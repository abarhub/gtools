package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/split"
)

var (
	sizeSplit       string
	bufferSizeSplit string
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "split file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := split.SplitParameters{File: args[0], SizeStr: sizeSplit,
			BufferSizeStr: bufferSizeSplit,
		}
		commandError = split.SplitCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureSplitCommandLine(rootCmd *cobra.Command) {
	splitCmd.Flags().StringVarP(&sizeSplit, "size", "s", "", "Size to split")
	splitCmd.Flags().StringVarP(&bufferSizeSplit, "bufferSize", "b", "", "Size of buffer")
	rootCmd.AddCommand(splitCmd)
}
