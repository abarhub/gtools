package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/unzip"
)

var (
	verboseUnzip     bool
	excludePathUnzip []string
	includePathUnzip []string
)

var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "unzip directory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := unzip.UnzipParameters{ZipFile: args[0], Directory: args[1],
			ExcludePath: excludePathUnzip, IncludePath: includePathUnzip, Verbose: verboseUnzip,
		}
		commandError = unzip.UnzipCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureUnzipCommandLine(rootCmd *cobra.Command) {
	unzipCmd.Flags().StringSliceVarP(&excludePathUnzip, "exclude", "e", []string{}, "Path to exclude")
	unzipCmd.Flags().StringSliceVarP(&includePathUnzip, "include", "i", []string{}, "Path to include")
	unzipCmd.Flags().BoolVarP(&verboseUnzip, "verbose", "v", false, "Verbose")
	rootCmd.AddCommand(unzipCmd)
}
