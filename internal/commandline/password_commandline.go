package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/password"
)

var (
	size     int
	number   bool
	letter   bool
	punctuer bool
	newLine  bool
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "generate password",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := password.PasswordParameters{Size: size, Number: number, Letter: letter,
			Punctuer: punctuer, NewLine: newLine}
		commandError = password.PasswordCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigurePasswordCommandLine(rootCmd *cobra.Command) {
	passwordCmd.Flags().IntVarP(&size, "size", "s", -1, "size of password")
	passwordCmd.Flags().BoolVarP(&number, "number", "n", true, "number in password")
	passwordCmd.Flags().BoolVarP(&letter, "letter", "l", true, "letter in password")
	passwordCmd.Flags().BoolVarP(&punctuer, "punctuer", "p", false, "punctuer in password")
	passwordCmd.Flags().BoolVarP(&punctuer, "newLine", "", true, "newline in password")
	rootCmd.AddCommand(passwordCmd)
}
