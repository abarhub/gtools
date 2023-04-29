package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
	"gtools/internal/command/rename"
)

func cmdRename(param rename.RenameParameters) error {
	return rename.RenameCommand(param)
}

func ConfigureRenameCommandLine(rootCmd *cobra.Command) {
	var (
		directory string
		verbose   bool
		dryRun    bool
		recusive  bool
	)

	var renameCmd = &cobra.Command{
		Use:   "rename [flags] files filesRenomed",
		Short: "rename files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			files := ""
			filesRenamed := ""
			if len(args) == 2 {
				files = args[0]
				filesRenamed = args[1]
			} else {
				return fmt.Errorf("rename need 2 arguments")
			}
			param := rename.RenameParameters{files, filesRenamed, recusive,
				verbose, dryRun, directory}
			commandError = cmdRename(param)
			return nil
		},
		DisableFlagsInUseLine: true,
	}

	renameCmd.Flags().StringVarP(&directory, "directory", "d", "",
		"directory source")
	renameCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose")
	renameCmd.Flags().BoolVarP(&dryRun, "dryRun", "y", false,
		"Don't rename")
	renameCmd.Flags().BoolVarP(&recusive, "recusive", "r", false,
		"sub directory")
	rootCmd.AddCommand(renameCmd)

}
