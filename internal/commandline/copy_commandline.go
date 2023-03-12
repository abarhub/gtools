package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
	copy2 "gtools/internal/command/copy"
)

var (
	excludePath    []string
	includePath    []string
	createDestDir  bool
	globDoubleStar bool
	verbose        bool
	dryRun         bool
	ignoreIfExist  string
)

const (
	paramExistsNo     = "no"
	paramExistsExists = "exists"
	paramExistsSize   = "size"
)

var copyCmd = &cobra.Command{
	Use:   "copy [flags] src dest",
	Short: "copy files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var fileExists copy2.FileExists
		switch ignoreIfExist {
		case paramExistsNo:
			fileExists = copy2.CopyFileExists
		case paramExistsExists:
			fileExists = copy2.NoCopyFileExists
		case paramExistsSize:
			fileExists = copy2.NoCopyFileExisteSizeFile
		default:
			return fmt.Errorf("Value '%v' for ignoreIfExists is invalide. Valid value is : %v, %v, %v",
				ignoreIfExist, paramExistsNo, paramExistsExists, paramExistsSize)

		}
		param := copy2.CopyParameters{args[0], args[1],
			excludePath, includePath, createDestDir,
			globDoubleStar, verbose, dryRun, fileExists}
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
	copyCmd.Flags().BoolVarP(&createDestDir, "createDestDir", "c", false,
		"Create destination directory if not exists")
	copyCmd.Flags().BoolVarP(&globDoubleStar, "doubleStar", "d", true,
		"Use global with double star for exclude and include")
	copyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show file copied")
	copyCmd.Flags().BoolVarP(&dryRun, "dryRun", "r", false, "Don't copy")
	copyCmd.Flags().StringVarP(&ignoreIfExist, "ignoreIfExists", "", paramExistsNo,
		"Copy if file exists. Value: "+paramExistsNo+"(always copy), "+paramExistsExists+"(if dest not exists), "+
			paramExistsSize+"(if size is différent)")
	rootCmd.AddCommand(copyCmd)

}
