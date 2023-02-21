package commandline

import (
	"github.com/spf13/cobra"
	copy2 "gtools/internal/command/copy"
)

var (
	excludePath    []string
	includePath    []string
	createDestDir  bool
	globDoubleStar bool
)

var copyCmd = &cobra.Command{
	Use:   "copy [flags] src dest",
	Short: "copy files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := copy2.CopyParameters{args[0], args[1],
			excludePath, includePath, createDestDir,
			globDoubleStar}
		commandError = cmdCopy(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func cmdCopy(param copy2.CopyParameters) error {
	return copy2.CopyDir(param)
}

func ConfigureCopyCommandLine(rootCmd *cobra.Command) {

	copyCmd.Flags().StringSliceVarP(&excludePath, "exclude", "e", []string{}, "Path to exclude")
	copyCmd.Flags().StringSliceVarP(&includePath, "include", "i", []string{}, "Path to include")
	copyCmd.Flags().BoolVarP(&createDestDir, "createDestDir", "c", false, "Create destination directory if not exists")
	copyCmd.Flags().BoolVarP(&globDoubleStar, "doubleStar", "d", true, "Use global with double star for exclude and include")
	rootCmd.AddCommand(copyCmd)

}
