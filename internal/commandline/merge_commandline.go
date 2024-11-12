package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/merge"
)

var (
	outputMerge string
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := merge.MergeParameters{File: args[0], OutputFile: outputMerge}
		commandError = merge.MergeCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureMergeCommandLine(rootCmd *cobra.Command) {
	mergeCmd.Flags().StringVarP(&outputMerge, "output", "o", "", "output file")
	rootCmd.AddCommand(mergeCmd)
}
