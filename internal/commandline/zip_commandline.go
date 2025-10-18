package commandline

import (
	"gtools/internal/command/zip"

	"github.com/spf13/cobra"
)

var (
	verboseZip     bool
	recursiveZip   bool
	excludePathZip []string
	includePathZip []string
)

var zipCmd = &cobra.Command{
	Use:   "zip [flags] zipfile file1 [file2 ...]",
	Short: "zip files",
	Long: `zip files
Example: zip -r archive.zip project/`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := zip.ZipParameters{ZipFile: args[0], Directory: args[1:],
			Recurvive: recursiveZip, ExcludePath: excludePathZip,
			IncludePath: includePathZip, Verbose: verboseZip,
		}
		commandError = zip.ZipCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureZipCommandLine(rootCmd *cobra.Command) {
	zipCmd.Flags().BoolVarP(&recursiveZip, "recursive", "r", true,
		"Zip subdirectory recursively")
	zipCmd.Flags().StringSliceVarP(&excludePathZip, "exclude", "e", []string{}, "Path to exclude")
	zipCmd.Flags().StringSliceVarP(&includePathZip, "include", "i", []string{}, "Path to include")
	zipCmd.Flags().BoolVarP(&verboseZip, "verbose", "v", false,
		"Verbose")
	rootCmd.AddCommand(zipCmd)
}
