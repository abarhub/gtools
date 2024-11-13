package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/aes"
)

var (
	inputFileAes   string
	outpoutFileAes string
	decryptAes     bool
	passwordAes    string
)

var aesCmd = &cobra.Command{
	Use:   "aes",
	Short: "aes encrypt/decrypt",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := aes.AesParameters{InputFile: inputFileAes, OutpoutFile: outpoutFileAes, Encrypt: !decryptAes,
			Password: passwordAes}
		commandError = aes.AesCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureAESCommandLine(rootCmd *cobra.Command) {
	aesCmd.Flags().StringVarP(&inputFileAes, "input", "i", "", "File input")
	aesCmd.Flags().StringVarP(&outpoutFileAes, "output", "o", "", "File output")
	aesCmd.Flags().BoolVarP(&decryptAes, "decrypt", "d", false, "decrypt. encrypt by default")
	aesCmd.Flags().StringVarP(&passwordAes, "password", "p", "", "Password for decrypt")
	rootCmd.AddCommand(aesCmd)
}
