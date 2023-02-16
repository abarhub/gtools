package commandline

import (
	"github.com/spf13/cobra"
	copy2 "gtools/internal/command/copy"
)

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

func cmdCopy(param copy2.CopyParameters) error {
	return copy2.CopyDir(param)
}

func ConfigureCopyCommandLine(rootCmd *cobra.Command) {

	copyCmd.Flags().StringArrayVarP(&excludePath, "exclude", "e", []string{}, "Path to exclude")
	copyCmd.Flags().StringArrayVarP(&includePath, "include", "i", []string{}, "Path to include")
	copyCmd.Flags().BoolVarP(&createDestDir, "createDestDir", "c", false, "Create destination directory if not exists")
	rootCmd.AddCommand(copyCmd)

}
